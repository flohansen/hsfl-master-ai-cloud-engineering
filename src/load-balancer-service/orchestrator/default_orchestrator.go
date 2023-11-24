package orchestrator

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"log"
	"math/rand"
	"net/url"
	"os"
	"reflect"
	"time"
)

const shutdownPollIntervalMax = 500 * time.Millisecond

type defaultOrchestrator struct {
	containers []string
}

func NewDefaultOrchestrator() *defaultOrchestrator {
	return &defaultOrchestrator{}
}

func (orchestrator *defaultOrchestrator) StartContainers(image string, replicas int, networkName string) []string {
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
		}, &container.HostConfig{PublishAllPorts: true}, &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				networkName: {},
			},
		}, nil, "")
		if err != nil {
			panic(err)
		}
		if err := cli.ContainerStart(context.Background(), createResponse.ID, types.ContainerStartOptions{}); err != nil {
			panic(err)
		}
		containers[i] = createResponse.ID
		log.Printf("Container created: %s", createResponse.ID)
	}

	orchestrator.containers = containers
	return containers
}

func (orchestrator *defaultOrchestrator) Shutdown(ctx context.Context) error {
	pollIntervalBase := time.Millisecond
	nextPollInterval := func() time.Duration {
		// Add 10% jitter.
		interval := pollIntervalBase + time.Duration(rand.Intn(int(pollIntervalBase/10)))
		// Double and clamp for next time.
		pollIntervalBase *= 2
		if pollIntervalBase > shutdownPollIntervalMax {
			pollIntervalBase = shutdownPollIntervalMax
		}
		return interval
	}

	timer := time.NewTimer(nextPollInterval())
	defer timer.Stop()
	for {
		if len(orchestrator.containers) > 0 {
			return orchestrator.StopAllContainers()
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			timer.Reset(nextPollInterval())
		}
	}

}

func (orchestrator *defaultOrchestrator) StopAllContainers() error {
	return orchestrator.StopContainers(orchestrator.containers...)
}

func (orchestrator *defaultOrchestrator) StopContainers(containers ...string) error {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}
	for _, currentContainer := range containers {
		if err := cli.ContainerRemove(context.Background(), currentContainer, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		}); err != nil {
			log.Printf(err.Error())
		}
		log.Printf("Container removed: %s", currentContainer)
	}
	return nil
}

func (orchestrator *defaultOrchestrator) GetContainerEndpoints(containers []string, networkName string) []*url.URL {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	endpoints := make([]*url.URL, len(containers))
	for i, c := range containers {
		inspectResponse, err := cli.ContainerInspect(context.Background(), c)
		if err != nil {
			panic(err)
		}
		var address = inspectResponse.NetworkSettings.Networks[networkName].IPAddress
		var port = reflect.ValueOf(inspectResponse.Config.ExposedPorts).MapKeys()[0].Interface().(nat.Port).Int()
		endpoint, err := url.Parse(fmt.Sprintf("http://%s:%d", address, port))
		if err != nil {
			panic(err)
		}

		endpoints[i] = endpoint
	}
	return endpoints
}
