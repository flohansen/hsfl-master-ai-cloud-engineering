package scheduler

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/balancer/endpoint"
)

type (
	// NewScheduler is a signature type for a Scheduler constructor
	NewScheduler func([]*endpoint.Endpoint) *Scheduler

	Scheduler interface {
		// SetEndpoints sets the reference to the endpoint list.
		//
		// It can be used for initialization or updating the reference.
		SetEndpoints(endpoints []*endpoint.Endpoint)
		// Next returns the next container endpoint to be used.
		//
		// Returns error if no endpoints are unavailable
		Next() (endpoint *endpoint.Endpoint, err error)
	}
)
