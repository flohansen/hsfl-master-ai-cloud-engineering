package balancer

import (
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type RandomBalancer struct {
	targets []http.Handler
}

func NewRandomBalancer(targetUrls []*url.URL) *RandomBalancer {
	targets := make([]http.Handler, len(targetUrls))

	for i, targetUrl := range targetUrls {
		targets[i] = httputil.NewSingleHostReverseProxy(targetUrl)
	}

	return &RandomBalancer{
		targets: targets,
	}
}

func (b *RandomBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	target := b.targets[rand.Intn(len(b.targets))]

	target.ServeHTTP(w, r)
}
