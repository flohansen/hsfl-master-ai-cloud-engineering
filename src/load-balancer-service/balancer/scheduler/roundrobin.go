package scheduler

import (
	"fmt"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/balancer/endpoint"
)

type RoundRobin struct {
	currentIndex int
	endpoints    []*endpoint.Endpoint
}

func (r RoundRobin) New(endpoints []*endpoint.Endpoint) Scheduler {
	var scheduler Scheduler = RoundRobin{
		currentIndex: -1,
		endpoints:    endpoints,
	}

	return scheduler
}

func (r RoundRobin) SetEndpoints(endpoints []*endpoint.Endpoint) {
	r.endpoints = endpoints
}

func (r RoundRobin) Next() (endpoint *endpoint.Endpoint, err error) {
	if len(r.endpoints) > 0 {
		defer func() {
			r.currentIndex = (r.currentIndex + 1) % len(r.endpoints)
		}()
		return r.endpoints[r.currentIndex], nil
	}
	return nil, fmt.Errorf("no endpoints are available")
}
