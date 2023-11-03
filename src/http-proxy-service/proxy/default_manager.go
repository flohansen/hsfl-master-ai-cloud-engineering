package proxy

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
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

type defaultManager struct {
	proxies []*httputil.ReverseProxy
	routing *router.Router
}

func NewDefaultManager(config *Config) *defaultManager {
	proxyManager := defaultManager{}

	// Prepare routing
	proxyManager.routing = router.New()

	for _, route := range config.ProxyRoutes {
		proxy, err := proxyManager.newProxy(route.Target)
		if err != nil {
			panic(err)
		}

		proxyManager.newHandler(proxy)
		proxyManager.proxies = append(proxyManager.proxies, proxy)

		proxyManager.routing.ALL(route.Context+"/*", proxyManager.newHandler(proxy))

		log.Printf("Mapping '%v' | %v ---> %v", route.Name, route.Context, route.Target)
	}

	return &proxyManager
}

func (dp defaultManager) GetProxyRouter() *router.Router {
	return dp.routing
}

func (dp defaultManager) newProxy(targetUrl string) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(targetUrl)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ModifyResponse = func(response *http.Response) error {
		dumpResponse, err := httputil.DumpResponse(response, false)
		if err != nil {
			return err
		}
		log.Println("Response: \r\n", string(dumpResponse))
		return nil
	}

	dp.proxies = append(dp.proxies, proxy)

	return proxy, nil
}

func (dp defaultManager) newHandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = r.Context().Value("wildcard0").(string)
		log.Println("Request URL: ", r.URL.String())
		p.ServeHTTP(w, r)
	}
}
