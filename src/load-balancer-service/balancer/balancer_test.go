package balancer

import (
	"github.com/stretchr/testify/assert"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/load-balancer-service/balancer/scheduler"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestBalancer(t *testing.T) {

	t.Run("Test NewBalancer", func(t *testing.T) {
		// Create a new balancer with a round-robin scheduler
		b := NewBalancer([]*url.URL{}, scheduler.NewRoundRobin)
		assert.NotNil(t, b)
		assert.NotNil(t, b.schedulerAlgorithm)
	})

	t.Run("Test ServeHTTP", func(t *testing.T) {
		// Create a new balancer with mock endpoints and a round-robin scheduler
		b := NewBalancer([]*url.URL{}, scheduler.NewRoundRobin)

		// Create a mock HTTP request
		req, err := http.NewRequest("GET", "/", nil)
		assert.NoError(t, err)

		// Create a mock HTTP response recorder
		w := httptest.NewRecorder()

		// Serve the request using the balancer
		b.ServeHTTP(w, req)

		// Assert that the status code is 503 (Service Unavailable) due to no available endpoints
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	})

	t.Run("Test SetHealthCheckFunction", func(t *testing.T) {
		// Create a new balancer with mock endpoints and a round-robin scheduler
		b := NewBalancer([]*url.URL{}, scheduler.NewRoundRobin)

		// Create a custom CheckFunction
		checkFunction := func(addr *url.URL) bool {
			return true
		}

		// Set the health check function for each endpoint in the balancer
		b.SetHealthCheckFunction(checkFunction, time.Second*2)

		// Assert that the health check function is set for each endpoint
		for _, ep := range b.endpoints {
			assert.NotNil(t, ep.HealthCheck)
			assert.Equal(t, checkFunction, ep.HealthCheck.IsAvailable())
		}
	})

	t.Run("Test ServeHTTP with Available Endpoint", func(t *testing.T) {
		mockURL, _ := url.Parse("http://example.com")
		b := NewBalancer([]*url.URL{mockURL}, scheduler.NewRoundRobin)

		// Create a mock HTTP request
		req, err := http.NewRequest("GET", "/", nil)
		assert.NoError(t, err)

		// Create a mock HTTP response recorder
		w := httptest.NewRecorder()

		// Serve the request using the balancer
		b.ServeHTTP(w, req)

		// Assert that the response status code is not 503 (Service Unavailable)
		assert.NotEqual(t, http.StatusServiceUnavailable, w.Code)
	})
}
