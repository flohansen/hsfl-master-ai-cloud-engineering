package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/load-balancer/balancer"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
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

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	err = CreateNetworkIfNotExists(cli, networkName)
	if err != nil {
		panic(err)
	}

	reader, err := cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	containerIds := StartContainers(cli, image, networkName, replicas)
	defer RemoveContainers(cli, containerIds)

	targetUrls, err := GetTargetUrls(cli, networkName, containerIds)
	if err != nil {
		panic(err)
	}

	loadBalancer := balancer.NewIPHashBalancer(targetUrls)

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

func StartContainers(cli *client.Client, imageName string, networkName string, replicas int) []string {
	containerConfig := &container.Config{
		Image: imageName,
	}

	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			networkName: {},
		},
	}

	containerIds := make([]string, replicas)

	for i := 0; i < replicas; i++ {
		resp, err := cli.ContainerCreate(context.Background(), containerConfig, nil, networkConfig, nil, "")
		if err != nil {
			log.Printf("Error creating container: %s\n", err)
		}

		containerIds[i] = resp.ID

		err = cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
		if err != nil {
			log.Printf("Error starting container: %s\n", err)
		}
	}

	return containerIds
}

func RemoveContainers(cli *client.Client, containerIDs []string) {
	log.Println("Stopping containers...")

	for _, containerID := range containerIDs {
		err := cli.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{
			Force: true,
		})
		if err != nil {
			log.Printf("Error stopping container %s: %s\n", containerID, err)
		}
	}
}

func CreateNetworkIfNotExists(cli *client.Client, networkName string) error {
	networks, err := cli.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		return err
	}

	networkExists := false
	for _, network := range networks {
		if network.Name == networkName {
			networkExists = true
			break
		}
	}

	if !networkExists {
		_, err = cli.NetworkCreate(context.Background(), networkName, types.NetworkCreate{})
		if err != nil {
			return err
		}
		log.Printf("Created network %s\n", networkName)
	}

	return nil
}

func GetTargetUrls(cli *client.Client, networkName string, containerIDs []string) ([]*url.URL, error) {
	targetUrls := make([]*url.URL, len(containerIDs))

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

		targetUrls[i] = targetUrl
	}

	return targetUrls, nil
}
