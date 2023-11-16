package random

import (
	"net/http"
	"testing"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/load-balancer/model"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// given
	r := New()

	// test
	assert.NotNil(t, r)
}

func TestGetTargetEmpty(t *testing.T) {
	// given
	r := New()

	// test
	r.GetTarget(&http.Request{}, []model.Target{}, func(_ model.Target) {
		t.Fatal("function should not be called with empty targets")
	})
}

func TestGetTargetNonEmpty(t *testing.T) {
	// given
	r := New()
	targets := []model.Target{{}, {}, {}}
	called := false

	// when
	r.GetTarget(&http.Request{}, targets, func(target model.Target) {
		called = true
		assert.Contains(t, targets, target)
	})

	// test
	assert.True(t, called, "callback function was not called")
}
