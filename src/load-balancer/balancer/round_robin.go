package balancer

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type RoundRobinBalancer struct {
	idx     int
	targets []http.Handler
}

func NewRoundRobinBalancer(targetUrls []*url.URL) *RoundRobinBalancer {
	targets := make([]http.Handler, len(targetUrls))

	for i, targetUrl := range targetUrls {
		targets[i] = httputil.NewSingleHostReverseProxy(targetUrl)
	}

	return &RoundRobinBalancer{
		targets: targets,
	}
}

func (b *RoundRobinBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	target := b.targets[b.idx]
	target.ServeHTTP(w, r)
	b.idx = (b.idx + 1) % len(b.targets)
}
