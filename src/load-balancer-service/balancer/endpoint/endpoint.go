package endpoint

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/balancer/endpoint/health"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Endpoint struct {
	url              *url.URL
	proxy            *httputil.ReverseProxy
	currentRequests  int
	lastResponseTime time.Duration
	healthCheck      *health.HealthCheck
}

func NewEndpoint(url *url.URL) *Endpoint {
	return &Endpoint{
		url:              url,
		proxy:            httputil.NewSingleHostReverseProxy(url),
		currentRequests:  0,
		lastResponseTime: 0,
	}
}

func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	e.currentRequests += 1
	defer func() {
		e.currentRequests -= 1
		e.lastResponseTime = time.Since(startTime)
	}()
	e.proxy.ServeHTTP(w, r)
}

func (e *Endpoint) IsAvailable() bool {
	if e.healthCheck != nil {
		return e.healthCheck.IsAvailable()
	}
	return true
}

func (e *Endpoint) SetHealthCheckFunction(check health.CheckFunction, period time.Duration) {
	e.healthCheck = health.NewHealthCheck(e.url, check, period)
}

func (e *Endpoint) GetCurrentRequests() int {
	return e.currentRequests
}

func (e *Endpoint) GetLastResponseTime() time.Duration {
	return e.lastResponseTime
}

func (e *Endpoint) GetURL() *url.URL {
	return e.url
}
