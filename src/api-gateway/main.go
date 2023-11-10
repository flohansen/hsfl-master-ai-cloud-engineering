package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/api-gateway/handler"
)

func main() {
	port := os.Getenv("PORT")

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

	// Start the Reverse Proxy
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
