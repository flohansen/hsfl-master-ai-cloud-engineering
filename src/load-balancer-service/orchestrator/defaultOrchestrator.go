package orchestrator

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"net/url"
	"os"
)

type defaultOrchestrator struct {
	containers []string
}

func NewDefaultOrchestrator() *defaultOrchestrator {
	return &defaultOrchestrator{}
}

func (orchestrator defaultOrchestrator) StartContainers(image string, replicas int) []string {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	pullResponse, err := cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	defer pullResponse.Close()

	io.Copy(os.Stdout, pullResponse)

	containers := make([]string, replicas)
	for i := 0; i < replicas; i++ {
		createResponse, err := cli.ContainerCreate(context.Background(), &container.Config{
			Image: image,
		}, &container.HostConfig{}, nil, nil, "")
		if err != nil {
			panic(err)
		}
		if err := cli.ContainerStart(context.Background(), createResponse.ID, types.ContainerStartOptions{}); err != nil {
			panic(err)
		}
		containers[i] = createResponse.ID
	}

	orchestrator.containers = containers
	return containers
}

func (orchestrator defaultOrchestrator) StopAllContainers() {
	orchestrator.StopContainers(orchestrator.containers)
}

func (orchestrator defaultOrchestrator) StopContainers(containers []string) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}
	for _, currentContainer := range containers {
		if err := cli.ContainerRemove(context.Background(), currentContainer, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			RemoveLinks:   true,
			Force:         true,
		}); err != nil {
			panic(err)
		}
	}
}

func (orchestrator defaultOrchestrator) GetContainerEndpoints(containers []string, networkName string) []*url.URL {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	endpoints := make([]*url.URL, len(containers))
	for i, container := range containers {
		inspectResponse, err := cli.ContainerInspect(context.Background(), container)
		if err != nil {
			panic(err)
		}
		//todo: adjust this for every orchestrator
		endpoint, err := url.Parse(fmt.Sprintf("http://%s:3000", inspectResponse.NetworkSettings.Networks[networkName].IPAddress))
		if err != nil {
			panic(err)
		}

		endpoints[i] = endpoint
	}
	return endpoints
}
