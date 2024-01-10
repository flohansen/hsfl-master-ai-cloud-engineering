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
		b := NewBalancer([]*url.URL{}, scheduler.NewRoundRobin)
		assert.NotNil(t, b)
		assert.NotNil(t, b.schedulerAlgorithm)
	})

	t.Run("Test ServeHTTP", func(t *testing.T) {
		b := NewBalancer([]*url.URL{}, scheduler.NewRoundRobin)

		req, err := http.NewRequest("GET", "/", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		b.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	})

	t.Run("Test SetHealthCheckFunction", func(t *testing.T) {
		b := NewBalancer([]*url.URL{}, scheduler.NewRoundRobin)

		checkFunction := func(addr *url.URL) bool {
			return true
		}

		b.SetHealthCheckFunction(checkFunction, time.Second*2)

		for _, ep := range b.endpoints {
			assert.NotNil(t, ep.HealthCheck)
			assert.Equal(t, checkFunction, ep.HealthCheck.IsAvailable())
		}
	})

	t.Run("Test ServeHTTP with Available Endpoint", func(t *testing.T) {
		mockURL, _ := url.Parse("http://example.com")
		b := NewBalancer([]*url.URL{mockURL}, scheduler.NewRoundRobin)
		req, err := http.NewRequest("GET", "/", nil)
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		b.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusServiceUnavailable, w.Code)
	})
}
