package main

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/http-proxy-service/proxy"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/http-proxy-service/proxy/proxyutils"
	"log"
	"net/http"
)

func main() {
	proxyConfig := proxyutils.ReadDefaultProxyManagerConfiguration("./config", "proxyConfig")
	proxyManager := proxy.NewDefaultManager(proxyConfig)

	log.Printf("Listening on %v", proxyConfig.ListenAddress)
	log.Fatal(http.ListenAndServe(proxyConfig.ListenAddress, proxyManager.GetProxyRouter()))
}
