package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/load-balancer/algorithm/round_robin.go"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/load-balancer/balancer"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/load-balancer/docker_helper"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/load-balancer/model"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func GetenvInt(key string) int {
	value := os.Getenv(key)
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}

	return valueInt
}

func main() {
	port := os.Getenv("PORT")
	image := os.Getenv("IMAGE")
	networkName := os.Getenv("NETWORK_NAME")
	replicas := GetenvInt("REPLICAS")
	healthCheckInterval := GetenvInt("HEALTH_CHECK_INTERVAL_SECONDS")

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	err = docker_helper.CreateNetworkIfNotExists(cli, networkName)
	if err != nil {
		panic(err)
	}

	reader, err := cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	containerIds := docker_helper.StartContainers(cli, image, networkName, replicas)
	defer docker_helper.RemoveContainers(cli, containerIds)

	targets, err := GetTargets(cli, networkName, containerIds)
	if err != nil {
		panic(err)
	}

	loadBalancer := balancer.NewBalancer(round_robin.New(), targets, healthCheckInterval)

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", port),
		Handler: loadBalancer,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		server.Shutdown(context.Background())
	}()

	server.ListenAndServe()
}

func GetTargets(cli *client.Client, networkName string, containerIDs []string) ([]model.Target, error) {
	targets := make([]model.Target, len(containerIDs))

	for i, containerID := range containerIDs {
		containerJson, err := cli.ContainerInspect(context.Background(), containerID)
		if err != nil {
			return nil, err
		}

		ip := containerJson.NetworkSettings.Networks[networkName].IPAddress

		urlStr := fmt.Sprintf("http://%s:%s", ip, "80")

		targetUrl, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}

		target := model.Target{
			ContainerId: containerID,
			Url:         targetUrl,
		}

		targets[i] = target
	}

	return targets, nil
}
