package main

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/http-proxy-service/proxy"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/http-proxy-service/proxy/proxyutils"
	"log"
	"net/http"
	"os"
)

var proxyConfigPath = "./config/proxyConfig.yaml"

func main() {
	log.Printf("Prepare service: http-proxy-service")

	var proxyConfig *proxy.Config

	envPath, isEnvPathSet := os.LookupEnv("PROXY_CONFIG_PATH")

	if isEnvPathSet {
		proxyConfigPath = envPath
	}

	proxyConfig = proxyutils.DefaultProxyManagerConfigurationReader(proxyConfigPath)

	log.Printf("Configuration loaded successfully with %d mappings", len(proxyConfig.ProxyRoutes))
	proxyManager := proxy.NewDefaultManager(proxyConfig)

	log.Printf("Listening on %v\r\n", proxyConfig.ListenAddress)
	log.Fatal(http.ListenAndServe(proxyConfig.ListenAddress, proxyManager.GetProxyRouter()))
}
