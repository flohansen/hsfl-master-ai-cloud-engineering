package endpoint

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Endpoint struct {
	proxy           *httputil.ReverseProxy
	currentRequests int
}

func NewEndpoint(url *url.URL) *Endpoint {
	return &Endpoint{
		proxy:           httputil.NewSingleHostReverseProxy(url),
		currentRequests: 0,
	}
}

func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.currentRequests += 1
	defer func() {
		e.currentRequests -= 1
	}()
	e.proxy.ServeHTTP(w, r)
}

func (e *Endpoint) IsAvailable() bool {
	return true
}

func (e *Endpoint) GetCurrentRequests() int {
	return e.currentRequests
}
