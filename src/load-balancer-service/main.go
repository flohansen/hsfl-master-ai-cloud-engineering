package main

import (
	"context"
	"flag"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/balancer"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/balancer/scheduler"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/orchestrator"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/utils"
	"net/http"
	"time"
)

func main() {
	defaultOrchestrator := orchestrator.NewDefaultOrchestrator()
	image := flag.String("image", "nginxdemos/hello", "")
	replicas := flag.Int("replicas", 5, "")
	networkName := flag.String("network", "bridge", "")
	flag.Parse()

	containers := defaultOrchestrator.StartContainers(*image, *replicas)

	loadBalancer := balancer.NewBalancer[scheduler.RoundRobin](defaultOrchestrator.GetContainerEndpoints(containers, *networkName))

	//Placeholder Server
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
