package scheduler

import (
	"github.com/stretchr/testify/assert"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/balancer/endpoint"
	"testing"
	"time"
)

func TestLeastResponseTimeScheduler(t *testing.T) {
	endpoint1 := &endpoint.Endpoint{}
	endpoint2 := &endpoint.Endpoint{}
	endpoints := []*endpoint.Endpoint{endpoint1, endpoint2}

	var testScheduler = NewLeastResponseTime(endpoints)

	t.Run(("Test SetEndpoints method"), func(t *testing.T) {
		newEndpoints := []*endpoint.Endpoint{endpoint1, endpoint2}
		(*testScheduler).SetEndpoints(endpoints)
		assert.Equal(t, newEndpoints, (*testScheduler).(*leastResponseTime).endpoints, "SetEndpoints method not working as expected")
	})

	t.Run(("Test next Method"), func(t *testing.T) {
		nextEndpoint, err := (*testScheduler).Next()
		assert.NoError(t, err, "Unexpected error in Next method")
		assert.NotNil(t, nextEndpoint, "Next method returned nil endpoint")
	})

	t.Run(("Test if next Method returns next Endpoint with least response time"), func(t *testing.T) {
		endpoint1 := &endpoint.Endpoint{CurrentRequests: 10, LastResponseTime: 2 * time.Second}
		endpoint2 := &endpoint.Endpoint{CurrentRequests: 5, LastResponseTime: 1 * time.Second}
		endpoints := []*endpoint.Endpoint{endpoint1, endpoint2}
		testScheduler = NewLeastResponseTime(endpoints)
		nextEndpoint, _ := (*testScheduler).Next()
		nextEndpoint2, _ := (*testScheduler).Next()
		assert.EqualValues(t, nextEndpoint, nextEndpoint2)
		assert.True(t, nextEndpoint2.LastResponseTime == 1*time.Second)
		endpoint2 = &endpoint.Endpoint{CurrentRequests: 5, LastResponseTime: 10 * time.Second}
		endpoints = []*endpoint.Endpoint{endpoint1, endpoint2}
		testScheduler = NewLeastResponseTime(endpoints)
		nextEndpoint2, _ = (*testScheduler).Next()
		assert.NotEqualValues(t, nextEndpoint, nextEndpoint2)
		assert.True(t, nextEndpoint2.LastResponseTime == 2*time.Second)
	})

	t.Run(("Test case where no available endpoints"), func(t *testing.T) {
		emptyScheduler := NewLeastResponseTime([]*endpoint.Endpoint{})
		emptyNextEndpoint, emptyErr := (*emptyScheduler).Next()
		assert.Error(t, emptyErr, "Expected error in Next method for empty testScheduler")
		assert.Nil(t, emptyNextEndpoint, "Expected nil endpoint for empty testScheduler")
	})
}
