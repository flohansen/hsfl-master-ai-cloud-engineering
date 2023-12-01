package docker_helper

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

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
