package ip_hash

import (
	"net/http"
	"testing"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/load-balancer/model"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// given
	ipHash := New()

	// test
	assert.NotNil(t, ipHash)
}

func TestGetTargetEmpty(t *testing.T) {
	// given
	ipHash := New()

	// test
	ipHash.GetTarget(&http.Request{}, []model.Target{}, func(_ model.Target) {
		t.Fatal("function should not be called with empty targets")
	})
}

func TestGetTargetNonEmpty(t *testing.T) {
	// given
	ipHash := New()
	targets := []model.Target{{}, {}, {}}

	// when
	r := &http.Request{RemoteAddr: "192.168.1.1:8080"}

	// test
	var firstTarget model.Target
	ipHash.GetTarget(r, targets, func(target model.Target) {
		firstTarget = target
	})
	ipHash.GetTarget(r, targets, func(target model.Target) {
		assert.Equal(t, firstTarget, target)
	})
}
