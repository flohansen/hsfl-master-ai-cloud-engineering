package least_connections

import (
	"net/http"
	"testing"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/load-balancer/model"
	"github.com/stretchr/testify/assert"
)

func TestLeastConnections(t *testing.T) {
	// given
	lc := New()

	targets := []model.Target{{ContainerId: "host1"}, {ContainerId: "host2"}, {ContainerId: "host3"}}

	lc.adjustConnectionCount(targets[0], 5)
	lc.adjustConnectionCount(targets[1], 3)
	lc.adjustConnectionCount(targets[2], 7)

	var targetHost string

	// when
	lc.GetTarget(&http.Request{}, targets, func(target model.Target) {
		targetHost = target.ContainerId
	})

	// test
	assert.Equal(t, "host2", targetHost)
}

func TestAdjustConnectionCount(t *testing.T) {
	// given
	lc := New()

	targets := []model.Target{{ContainerId: "host1"}, {ContainerId: "host2"}, {ContainerId: "host3"}}

	// when
	lc.adjustConnectionCount(targets[0], 1)
	lc.adjustConnectionCount(targets[1], -1)
	lc.adjustConnectionCount(targets[2], 3)

	// test
	assert.Equal(t, 1, lc.count[targets[0]])
	assert.Equal(t, -1, lc.count[targets[1]])
	assert.Equal(t, 3, lc.count[targets[2]])
}
