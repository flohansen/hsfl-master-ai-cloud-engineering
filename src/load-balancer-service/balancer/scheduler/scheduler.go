package scheduler

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/balancer/endpoint"
)

type Scheduler interface {
	// New initializes a new instance of a scheduler and sets the reference to the endpoints list
	New(endpoints []*endpoint.Endpoint) Scheduler
	// SetEndpoints sets the reference to the endpoint list.
	//
	// It can be used for initialization or updating the reference.
	SetEndpoints(endpoints []*endpoint.Endpoint)
	// Next returns the next container endpoint to be used.
	//
	// Returns error if no endpoints are unavailable
	Next() (endpoint *endpoint.Endpoint, err error)
}
