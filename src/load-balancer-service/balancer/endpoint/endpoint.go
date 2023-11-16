package endpoint

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Endpoint struct {
	proxy            *httputil.ReverseProxy
	currentRequests  int
	lastResponseTime time.Duration
}

func NewEndpoint(url *url.URL) *Endpoint {
	return &Endpoint{
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
	return true
}

func (e *Endpoint) GetCurrentRequests() int {
	return e.currentRequests
}

func (e *Endpoint) GetLastResponseTime() time.Duration {
	return e.lastResponseTime
}
