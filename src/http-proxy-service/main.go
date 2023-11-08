package main

import (
	"github.com/spf13/viper"
	router "hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ProxyRoute struct {
	Name    string `mapstructure:"name"`
	Context string `mapstructure:"context"`
	Target  string `mapstructure:"target"`
}

type ProxyConfig struct {
	ListenAddress string       `mapstructure:"listenAddress"`
	ProxyRoutes   []ProxyRoute `mapstructure:"proxyRoutes"`
}

func main() {
	// Load config
	proxyConfig := readProxyConfiguration()

	// Prepare routing
	routing := router.New()

	for _, route := range proxyConfig.ProxyRoutes {
		proxy, err := newProxy(route.Target)
		if err != nil {
			panic(err)
		}

		log.Printf("Mapping '%v' | %v ---> %v", route.Name, route.Context, route.Target)

		routing.ALL(route.Context+"/*", newHandler(proxy))
	}

	log.Printf("Listening on %v", proxyConfig.ListenAddress)
	log.Fatal(http.ListenAndServe(proxyConfig.ListenAddress, routing))
}

func readProxyConfiguration() *ProxyConfig {
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")
	viper.SetConfigName("proxyConfig")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while loading proxy configuration: %v", err)
	}
	viper.AutomaticEnv()

	proxyConfig := &ProxyConfig{}

	err = viper.UnmarshalKey("proxy", proxyConfig)
	if err != nil {
		panic(err)
	}

	return proxyConfig
}

func newProxy(targetUrl string) (*httputil.ReverseProxy, error) {
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
	return proxy, nil
}

func newHandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = r.Context().Value("wildcard0").(string)
		log.Println("Request URL: ", r.URL.String())
		p.ServeHTTP(w, r)
	}
}
