package orchestrator

import "net/url"

type Orchestrator interface {
	StartContainers(image string, replicas int) []string
	StopContainers(containers []string)
	StopAllContainers()
	GetContainerEndpoints(containers []string, networkName string) []*url.URL
}
