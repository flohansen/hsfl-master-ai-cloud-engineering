package proxy

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Route struct {
	Name    string `mapstructure:"name"`
	Context string `mapstructure:"context"`
	Target  string `mapstructure:"target"`
}

type Config struct {
	ListenAddress string  `mapstructure:"listenAddress"`
	ProxyRoutes   []Route `mapstructure:"proxyRoutes"`
}

type DefaultManager struct {
	proxies []*httputil.ReverseProxy
	routing *router.Router
}

func NewDefaultManager(config *Config) *DefaultManager {
	proxyManager := DefaultManager{}

	// Prepare routing
	proxyManager.routing = router.New()

	for _, route := range config.ProxyRoutes {
		proxy, err := proxyManager.newProxy(route.Target)
		if err != nil {
			panic(err)
		}

		proxyManager.newHandler(proxy)
		proxyManager.proxies = append(proxyManager.proxies, proxy)

		if !strings.HasSuffix(route.Context, "/") {
			route.Context = route.Context + "/"
		}

		proxyManager.routing.ALL(route.Context+"*", proxyManager.newHandler(proxy))

		log.Printf("Mapping '%v' | %v ---> %v", route.Name, route.Context, route.Target)
	}

	return &proxyManager
}

func (dp DefaultManager) GetProxyRouter() *router.Router {
	return dp.routing
}

func (dp DefaultManager) newProxy(targetUrl string) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(targetUrl)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ModifyResponse = func(response *http.Response) error {
		response.Header.Set("X-Proxy", "Price Whisper Proxy")
		dumpResponse, err := httputil.DumpResponse(response, false)
		if err != nil {
			return err
		}
		log.Println("Relaying to: ", target)
		log.Println("Response to client: \r\n", string(dumpResponse))
		return nil
	}

	dp.proxies = append(dp.proxies, proxy)

	return proxy, nil
}

func (dp DefaultManager) newHandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request URL: %v ---> /%v\n", r.URL.String(), r.Context().Value("wildcard0").(string))
		r.Header.Set("X-Proxy", "Price Whisper Proxy")
		r.Header.Set("X-Forwarded-For", strings.Split(r.RemoteAddr, ":")[0])
		r.Header.Set("X-Forwarded-Host", r.Host)
		r.URL.Path = r.Context().Value("wildcard0").(string)
		p.ServeHTTP(w, r)
	}
}
