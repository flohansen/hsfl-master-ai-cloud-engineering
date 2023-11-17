package scheduler

import (
	"fmt"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/balancer/endpoint"
	"sort"
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
		sortedEndpoints := make([]*endpoint.Endpoint, len(r.endpoints))
		copy(sortedEndpoints, r.endpoints)

		sort.Slice(sortedEndpoints, func(i, j int) bool {
			return sortedEndpoints[i].GetLastResponseTime() < sortedEndpoints[j].GetLastResponseTime()
		})

		for _, ep := range sortedEndpoints {
			if ep.IsAvailable() {
				return ep, nil
			}
		}
	}

	return nil, fmt.Errorf("no available endpoints")
}
