package round_robin

import (
	"net/http"
	"testing"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/load-balancer/model"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// given
	rr := New()

	// test
	assert.Equal(t, 0, rr.idx)
}

func TestGetTargetEmpty(t *testing.T) {
	// given
	rr := New()

	// test
	rr.GetTarget(&http.Request{}, []model.Target{}, func(_ model.Target) {
		t.Fatal("function should not be called with empty targets")
	})
}

func TestGetTargetNonEmpty(t *testing.T) {
	// given
	rr := New()
	targets := []model.Target{{}, {}, {}}

	// test
	rr.GetTarget(&http.Request{}, targets, func(target model.Target) {
		assert.Equal(t, targets[0], target)
	})
	rr.GetTarget(&http.Request{}, targets, func(target model.Target) {
		assert.Equal(t, targets[1], target)
	})
	rr.GetTarget(&http.Request{}, targets, func(target model.Target) {
		assert.Equal(t, targets[2], target)
	})
}

func TestGetTargetWrapAround(t *testing.T) {
	// given
	rr := New()
	targets := []model.Target{{}, {}}

	// test
	rr.GetTarget(&http.Request{}, targets, func(target model.Target) {
		assert.Equal(t, targets[0], target)
	})
	rr.GetTarget(&http.Request{}, targets, func(target model.Target) {
		assert.Equal(t, targets[1], target)
	})
	rr.GetTarget(&http.Request{}, targets, func(target model.Target) {
		assert.Equal(t, targets[0], target)
	})
}
