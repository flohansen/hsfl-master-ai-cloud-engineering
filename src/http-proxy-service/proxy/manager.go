package proxy

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router"
	"net/http"
	"net/http/httputil"
)

type Manager interface {
	GetProxyRouter() *router.Router
	newProxy(targetUrl string) (*httputil.ReverseProxy, error)
	newHandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request)
}
