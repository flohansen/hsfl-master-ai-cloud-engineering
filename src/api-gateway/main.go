package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/api-gateway/handler"
)

func main() {
	port := os.Getenv("SERVER_PORT")

	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	bulletinBoardServiceURL := os.Getenv("BULLETIN_BOARD_SERVICE_URL")
	feedServiceURL := os.Getenv("FEED_SERVICE_URL")
	webServiceURL := os.Getenv("WEB_SERVICE_URL")

	// Configuration for the Reverse Proxy
	config := handler.ReverseProxyConfig{
		Services: []handler.Service{
			{Name: "frontend", ContextPath: "/", TargetURL: webServiceURL},
			{Name: "auth", ContextPath: "/auth", TargetURL: authServiceURL},
			{Name: "bulletin-board", ContextPath: "/bulletin-board", TargetURL: bulletinBoardServiceURL},
			{Name: "feed", ContextPath: "/feed", TargetURL: feedServiceURL},
		},
	}

	// Create a new Reverse Proxy
	reverseProxy := handler.NewReverseProxy(config)

	// Use the Reverse Proxy as a http.Handler in http.ListenAndServe
	http.Handle("/", reverseProxy)

	// Start the Reverse Proxy
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
