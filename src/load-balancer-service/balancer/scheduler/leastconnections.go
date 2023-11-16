package scheduler

import (
	"fmt"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/balancer/endpoint"
	"math"
)

type leastConnections struct {
	requestsCount map[string]uint
	endpoints     []*endpoint.Endpoint
}

func NewLeastConnections(endpoints []*endpoint.Endpoint) *Scheduler {
	var scheduler Scheduler = &leastConnections{
		requestsCount: make(map[string]uint),
		endpoints:     endpoints,
	}
	return &scheduler
}

func (r *leastConnections) SetEndpoints(endpoints []*endpoint.Endpoint) {
	r.endpoints = endpoints
}

func (r *leastConnections) Next() (*endpoint.Endpoint, error) {
	if len(r.endpoints) > 0 {
		minRequests := math.MaxUint32
		var leastConnEndpoint *endpoint.Endpoint

		for _, ep := range r.endpoints {
			currentRequests := ep.GetCurrentRequests()

			// Update the least connection endpoint if the current one has fewer requests
			if currentRequests < minRequests {
				minRequests = currentRequests
				leastConnEndpoint = ep
			}
		}

		if leastConnEndpoint != nil {
			return leastConnEndpoint, nil
		}
	}

	return nil, fmt.Errorf("no endpoints are available")
}
