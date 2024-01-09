package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/api-gateway/config"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/api-gateway/handler"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cfg := config.Config{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("error while parsing enviroment variables: %s", err.Error())
	}

	// Configuration for the Reverse Proxy
	config := handler.ReverseProxyConfig{
		Services: []handler.Service{
			{Name: "frontend", ContextPath: "/", TargetURL: cfg.WebServiceUrl},
			{Name: "auth", ContextPath: "/auth", TargetURL: cfg.AuthServiceUrl},
			{Name: "bulletin-board", ContextPath: "/bulletin-board", TargetURL: cfg.BulletinBoardServiceUrl},
			{Name: "feed", ContextPath: "/feed", TargetURL: cfg.FeedServiceUrl},
		},
	}

	// Create a new Reverse Proxy
	reverseProxy := handler.NewReverseProxy(config)

	// Use the Reverse Proxy as a http.Handler in http.ListenAndServe
	http.Handle("/", reverseProxy)

	// Start the Reverse Proxy
	log.Printf("Starting HTTP server on port %s", cfg.HttpServerPort)
	addr := fmt.Sprintf("0.0.0.0:%s", cfg.HttpServerPort)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
