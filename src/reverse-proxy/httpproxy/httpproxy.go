package httpproxy

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

type RouteMapping struct {
	hostIndex int
	path      *regexp.Regexp
	hosts     []*url.URL
}

type HTTPProxy struct {
	Mappings []*RouteMapping
}

func New() *HTTPProxy {
	return &HTTPProxy{}
}

func (p *HTTPProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, mapping := range p.Mappings {
		matches := mapping.path.FindAllStringSubmatch(r.URL.Path, -1)
		if len(matches) > 0 {
			host := mapping.hosts[mapping.hostIndex]

			r.Header.Set("X-Forwarded-For", r.RemoteAddr)
			r.Header.Set("X-Forwarded-Host", r.URL.Host)
			r.Host = host.Host
			r.URL.Host = host.Host
			r.URL.Scheme = host.Scheme
			r.RequestURI = ""

			originServerResponse, err := http.DefaultClient.Do(r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = fmt.Fprint(w, err)
				mapping.hostIndex = (mapping.hostIndex + 1) % len(mapping.hosts)
				return
			}
			w.WriteHeader(http.StatusOK)
			io.Copy(w, originServerResponse.Body)

			// Round-robin
			mapping.hostIndex = (mapping.hostIndex + 1) % len(mapping.hosts)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	return
}

func (p *HTTPProxy) Add(path string, hosts []string) error {
	pattern := regexp.MustCompile(path)

	if len(hosts) < 1 {
		return errors.New("there was no host provided")
	}

	var urls []*url.URL
	for _, hostAddr := range hosts {
		host, err := url.Parse(hostAddr)
		if err != nil {
			return errors.New("invalid origin server URL")
		}
		urls = append(urls, host)
	}

	p.Mappings = append(p.Mappings, &RouteMapping{0, pattern, urls})
	return nil
}
