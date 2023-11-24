package orchestrator

import (
	"context"
	"net/url"
)

type Orchestrator interface {
	StartContainers(image string, replicas int, networkName string) []string
	StopContainers(containers ...string) error
	StopAllContainers() error
	Shutdown(ctx context.Context) error
	GetContainerEndpoints(containers []string, networkName string) []*url.URL
}
