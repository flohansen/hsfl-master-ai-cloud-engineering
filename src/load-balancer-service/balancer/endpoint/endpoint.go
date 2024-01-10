package endpoint

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/balancer/endpoint/health"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Endpoint struct {
	Url              *url.URL
	Proxy            *httputil.ReverseProxy
	CurrentRequests  int
	LastResponseTime time.Duration
	HealthCheck      *health.HealthCheck
}

func NewEndpoint(url *url.URL) *Endpoint {
	return &Endpoint{
		Url:              url,
		Proxy:            httputil.NewSingleHostReverseProxy(url),
		CurrentRequests:  0,
		LastResponseTime: 0,
	}
}

func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	e.CurrentRequests += 1
	defer func() {
		e.CurrentRequests -= 1
		e.LastResponseTime = time.Since(startTime)
	}()
	e.Proxy.ServeHTTP(w, r)
}

func (e *Endpoint) IsAvailable() bool {
	if e.HealthCheck != nil {
		return e.HealthCheck.IsAvailable()
	}
	return true
}

func (e *Endpoint) SetHealthCheckFunction(check health.CheckFunction, period time.Duration) {
	e.HealthCheck = health.NewHealthCheck(e.Url, check, period)
}

func (e *Endpoint) GetCurrentRequests() int {
	return e.CurrentRequests
}

func (e *Endpoint) GetLastResponseTime() time.Duration {
	return e.LastResponseTime
}

func (e *Endpoint) GetURL() *url.URL {
	return e.Url
}
