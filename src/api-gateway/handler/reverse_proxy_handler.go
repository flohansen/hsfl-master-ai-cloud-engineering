package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
)

type Service struct {
	Name         string
	ContextPath  string
	TargetURL    string
	ReverseProxy *httputil.ReverseProxy
}

// ReverseProxyConfig represents the configuration for the Reverse Proxy.
type ReverseProxyConfig struct {
	Services []Service
}

// ReverseProxyHandler is a struct representing the Reverse Proxy.
type ReverseProxyHandler struct {
	Services map[*regexp.Regexp]*httputil.ReverseProxy
}

func NewReverseProxy(config ReverseProxyConfig) *ReverseProxyHandler {
	services := make(map[*regexp.Regexp]*httputil.ReverseProxy)

	for _, service := range config.Services {
		proxy := newSingleHostReverseProxy(service.TargetURL)
		pattern, err := regexp.Compile(service.ContextPath)

		if err != nil {
			panic(err)
		}

		services[pattern] = proxy
	}

	return &ReverseProxyHandler{
		Services: services,
	}
}
func newSingleHostReverseProxy(targetURL string) *httputil.ReverseProxy {
	target, _ := url.Parse(targetURL)
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
		http.Error(writer, "Proxy error", http.StatusInternalServerError)
		log.Printf("proxy for %s error: %v", targetURL, e)
	}
	return proxy
}

func (rp *ReverseProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for pattern, proxy := range rp.Services {
		if pattern.MatchString(r.URL.Path) {
			log.Printf("Proxying %s to %s", r.URL.Path, pattern.String())
			proxy.ServeHTTP(w, r)
			return
		}
	}

	http.Error(w, "Service not found", http.StatusNotFound)
}
