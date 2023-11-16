package scheduler

import (
	"fmt"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/balancer/endpoint"
)

type leastResponseTime struct {
	endpoints []*endpoint.Endpoint
}

func NewLeastResponseTime(endpoints []*endpoint.Endpoint) *Scheduler {
	var scheduler Scheduler = &leastResponseTime{
		endpoints: endpoints,
	}
	return &scheduler
}

func (r *leastResponseTime) SetEndpoints(endpoints []*endpoint.Endpoint) {
	r.endpoints = endpoints
}

func (r *leastResponseTime) Next() (*endpoint.Endpoint, error) {
	if len(r.endpoints) > 0 {
		var leastResponseEndpoint *endpoint.Endpoint
		var minResponseTime = r.endpoints[0].GetLastResponseTime()

		for _, ep := range r.endpoints {
			currentResponseTime := ep.GetLastResponseTime()

			if currentResponseTime < minResponseTime {
				minResponseTime = currentResponseTime
				leastResponseEndpoint = ep
			}
		}

		if leastResponseEndpoint != nil {
			return leastResponseEndpoint, nil
		}
	}

	return nil, fmt.Errorf("no endpoints are available")
}
