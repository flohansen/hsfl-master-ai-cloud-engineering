package main

import (
	"fmt"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/api-gateway/handler"
	"net/http"
)

func main() {
	// Configuration for the Reverse Proxy
	config := handler.ReverseProxyConfig{
		Services: []handler.Service{
			{Name: "auth", ContextPath: "/auth", TargetURL: "http://localhost:8081"},
			{Name: "bulletinboard", ContextPath: "/bulletin-board", TargetURL: "http://localhost:8080"},
		},
	}

	// Create a new Reverse Proxy
	reverseProxy := handler.NewReverseProxy(config)

	// Use the Reverse Proxy as an http.Handler in http.ListenAndServe
	http.Handle("/", reverseProxy)

	// Start the Reverse Proxy on a port
	port := 8000
	fmt.Printf("Reverse Proxy is running on port %d\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
