package main

import (
	"context"
	"flag"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/balancer"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/balancer/scheduler"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/orchestrator"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	defaultOrchestrator := orchestrator.NewDefaultOrchestrator()
	image := flag.String("image", "hello-world", "")
	replicas := flag.Int("replicas", 1, "")
	networkName := flag.String("network", "bridge", "")
	flag.Parse()

	containers := defaultOrchestrator.StartContainers(*image, *replicas)

	loadBalancer := balancer.NewBalancer[scheduler.RoundRobin](defaultOrchestrator.GetContainerEndpoints(containers, *networkName))

	//Placeholder Server
	server := &http.Server{
		Addr:    ":3000",
		Handler: loadBalancer,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		server.Shutdown(context.Background())
		defaultOrchestrator.StopContainers(containers)
	}()

	log.Fatal(server.ListenAndServe())
}
