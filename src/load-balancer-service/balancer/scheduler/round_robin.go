package scheduler

import (
	"fmt"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/balancer/endpoint"
)

type roundRobin struct {
	currentIndex int
	endpoints    []*endpoint.Endpoint
}

func NewRoundRobin(endpoints []*endpoint.Endpoint) *Scheduler {
	var scheduler Scheduler = &roundRobin{
		currentIndex: -1,
		endpoints:    endpoints,
	}

	return &scheduler
}

func (r *roundRobin) SetEndpoints(endpoints []*endpoint.Endpoint) {
	r.endpoints = endpoints
}

func (r *roundRobin) Next() (endpoint *endpoint.Endpoint, err error) {
	if len(r.endpoints) > 0 {
		stopIndex := r.currentIndex
		r.currentIndex = (r.currentIndex + 1) % len(r.endpoints)
		for i := r.currentIndex; !r.endpoints[i].IsAvailable() && r.currentIndex != stopIndex; i++ {
			r.currentIndex = (r.currentIndex + i) % len(r.endpoints)
		}

		if r.endpoints[r.currentIndex].IsAvailable() {
			return r.endpoints[r.currentIndex], nil
		}
	}
	return nil, fmt.Errorf("no endpoints are available")
}
