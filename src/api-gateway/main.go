package main

import (
	"fmt"
	"net/http"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/api-gateway/handler"

)


func main() {
	// Configuration for the Reverse Proxy
	config := handler.ReverseProxyConfig{

		//TODO CORRECT PORTS AND PATHS

		//AuthServiceURL:       "http://authservice:port",
		//BulletinBoardServiceURL: "http://bulletinboardservice:port",
		FeedServiceURL:         "http://localhost:8081",
	}

	// Create a new Reverse Proxy
	reverseProxy := handler.NewReverseProxy(config)

	// Use the Reverse Proxy as an http.Handler in http.ListenAndServe
	http.Handle("/", reverseProxy)

	// Start the Reverse Proxy on a port
	port := 8080
	fmt.Printf("Reverse Proxy is running on port %d\n", port)
	
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}