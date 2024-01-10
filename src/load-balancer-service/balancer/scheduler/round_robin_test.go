package scheduler

import (
	"github.com/stretchr/testify/assert"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/balancer/endpoint"
	"testing"
)

func TestRoundRobinScheduler(t *testing.T) {

	endpoint1 := &endpoint.Endpoint{}
	endpoint2 := &endpoint.Endpoint{}
	endpoint3 := &endpoint.Endpoint{}
	endpoints := []*endpoint.Endpoint{endpoint1, endpoint2, endpoint3}

	var testScheduler = NewRoundRobin(endpoints)

	t.Run(("Test SetEndpoints method"), func(t *testing.T) {
		newEndpoints := []*endpoint.Endpoint{endpoint1, endpoint2, endpoint3}
		(*testScheduler).SetEndpoints(endpoints)
		assert.Equal(t, newEndpoints, (*testScheduler).(*roundRobin).endpoints, "SetEndpoints method not working as expected")
	})

	t.Run(("Test next Method"), func(t *testing.T) {
		nextEndpoint, err := (*testScheduler).Next()
		assert.NoError(t, err, "Unexpected error in Next method")
		assert.NotNil(t, nextEndpoint, "Next method returned nil endpoint")
	})

	t.Run(("Test if next Method returns next Endpoint"), func(t *testing.T) {
		nextEndpoint, _ := (*testScheduler).Next()
		nextEndpoint2, _ := (*testScheduler).Next()
		assert.EqualValues(t, nextEndpoint, nextEndpoint2)
	})

	t.Run(("Test case where no available endpoints"), func(t *testing.T) {
		emptyScheduler := NewLeastResponseTime([]*endpoint.Endpoint{})
		emptyNextEndpoint, emptyErr := (*emptyScheduler).Next()
		assert.Error(t, emptyErr, "Expected error in Next method for empty testScheduler")
		assert.Nil(t, emptyNextEndpoint, "Expected nil endpoint for empty testScheduler")
	})
}
