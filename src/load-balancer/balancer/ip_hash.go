package balancer

import (
	"hash/fnv"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
)

type IPHashBalancer struct {
	targets []http.Handler
}

func NewIPHashBalancer(targetUrls []*url.URL) *IPHashBalancer {
	targets := make([]http.Handler, len(targetUrls))

	for i, targetUrl := range targetUrls {
		targets[i] = httputil.NewSingleHostReverseProxy(targetUrl)
	}

	return &IPHashBalancer{targets: targets}
}

func (b *IPHashBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`).FindString(r.RemoteAddr)

	if ip == "" {
		http.Error(w, "Could not find IP address", http.StatusInternalServerError)
		return
	}

	hashValue := hash(ip)
	target := b.targets[hashValue%uint32(len(b.targets))]

	target.ServeHTTP(w, r)
}

func hash(str string) uint32 {
	hasher := fnv.New32a()
	hasher.Write([]byte(str))

	return hasher.Sum32()
}
