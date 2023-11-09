package proxy

import (
	router "hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib"
	"net/http"
	"net/http/httputil"
)

type Manager interface {
	GetProxyRouter() *router.Router
	newProxy(targetUrl string) (*httputil.ReverseProxy, error)
	newHandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request)
}
