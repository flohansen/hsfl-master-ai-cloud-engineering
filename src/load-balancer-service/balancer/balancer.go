package balancer

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/balancer/endpoint"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/balancer/endpoint/health"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/balancer/scheduler"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Balancer struct {
	endpoints          []*endpoint.Endpoint
	schedulerAlgorithm *scheduler.Scheduler
}

func NewBalancer(targetUrls []*url.URL, scheduler scheduler.NewScheduler) *Balancer {
	endpoints := make([]*endpoint.Endpoint, len(targetUrls))
	for i, targetUrl := range targetUrls {
		endpoints[i] = endpoint.NewEndpoint(targetUrl)
	}

	return &Balancer{
		endpoints:          endpoints,
		schedulerAlgorithm: scheduler(endpoints),
	}
}

func (b *Balancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	endpoint, err := (*b.schedulerAlgorithm).Next()
	if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	endpoint.ServeHTTP(w, r)
}

func (b *Balancer) SetHealthCheckFunction(check health.CheckFunction, period time.Duration) {
	for _, endpoint := range b.endpoints {
		endpoint.SetHealthCheckFunction(check, period)
	}
}
