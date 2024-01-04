package main

import (
	"context"
	"flag"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/balancer"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/balancer/endpoint/health"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/balancer/scheduler"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/orchestrator"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/utils"
	"log"
	"net/http"
	"time"
)

func main() {
	// Configuration
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()
	configuration := utils.GetConfig(*configPath)

	// Log configuration
	log.Printf("Image is: %s", configuration.Image)
	log.Printf("No. of Replicas is: %d", configuration.Replicas)
	log.Printf("Network is: %s", configuration.Network)

	// Initialize orchestrator and start containers
	var defaultOrchestrator orchestrator.Orchestrator = orchestrator.NewDefaultOrchestrator()
	containers := defaultOrchestrator.StartContainers(configuration.Image, configuration.Replicas, configuration.Network)

	// Initialize load balancer
	loadBalancer := balancer.NewBalancer(
		defaultOrchestrator.GetContainerEndpoints(containers, configuration.Network),
		scheduler.NewRoundRobin)
	loadBalancer.SetHealthCheckFunction(health.DefaultHealthCheck, 5*time.Second)

	// Start web server
	server := &http.Server{
		Addr:    ":3000",
		Handler: loadBalancer,
	}

	//
	go func() {
		log.Print(server.ListenAndServe())
	}()

	// Graceful shutdown
	wait := utils.GracefulShutdown(context.Background(), 30*time.Second, map[string]utils.Operation{
		"orchestrator": func(ctx context.Context) error {
			return defaultOrchestrator.Shutdown(ctx)
		},
		"load-balancer": func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})

	<-wait
}
