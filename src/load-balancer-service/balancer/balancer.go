package balancer

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/balancer/endpoint"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/balancer/scheduler"
	"net/http"
	"net/url"
)

type Balancer struct {
	endpoints          []*endpoint.Endpoint
	schedulerAlgorithm scheduler.Scheduler
}

func NewBalancer[T scheduler.Scheduler](targetUrls []*url.URL) *Balancer {
	endpoints := make([]*endpoint.Endpoint, len(targetUrls))
	for i, targetUrl := range targetUrls {
		endpoints[i] = endpoint.NewEndpoint(targetUrl)
	}

	var schedulerAlgorithm T

	return &Balancer{
		endpoints:          endpoints,
		schedulerAlgorithm: T.New(schedulerAlgorithm, endpoints),
	}
}

func (b *Balancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	endpoint, err := b.schedulerAlgorithm.Next()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	endpoint.ServeHTTP(w, r)
}
