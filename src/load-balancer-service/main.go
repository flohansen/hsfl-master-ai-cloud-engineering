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
	defaultOrchestrator := orchestrator.NewDefaultOrchestrator()
	image := flag.String("image", "nginxdemos/hello", "")
	replicas := flag.Int("replicas", 5, "")
	networkName := flag.String("network", "bridge", "")
	flag.Parse()
	log.Printf("Image is: %s", *image)
	log.Printf("No. of Replicas is: %d", *replicas)
	log.Printf("Network is: %s", *networkName)

	containers := defaultOrchestrator.StartContainers(*image, *replicas, *networkName)

	loadBalancer := balancer.NewBalancer(
		defaultOrchestrator.GetContainerEndpoints(containers, *networkName),
		scheduler.NewLeastResponseTime)
	loadBalancer.SetHealthCheckFunction(health.DefaultHealthCheck, 5*time.Second)

	server := &http.Server{
		Addr:    ":3000",
		Handler: loadBalancer,
	}

	go func() {
		log.Print(server.ListenAndServe())
	}()

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
