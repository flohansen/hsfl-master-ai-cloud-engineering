package main

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/http-proxy-service/proxy"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/http-proxy-service/proxy/proxyutils"
	"log"
	"net/http"
)

func main() {
	log.Printf("Prepare service: http-proxy-service")

	proxyConfig := proxyutils.DefaultProxyManagerConfigurationReader("./config", "proxyConfig.yaml")
	log.Printf("Configuration loaded successfully with %d mappings", len(proxyConfig.ProxyRoutes))
	proxyManager := proxy.NewDefaultManager(proxyConfig)

	log.Printf("Listening on %v\r\n", proxyConfig.ListenAddress)
	log.Fatal(http.ListenAndServe(proxyConfig.ListenAddress, proxyManager.GetProxyRouter()))
}
