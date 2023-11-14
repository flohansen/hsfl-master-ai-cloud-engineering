package handler

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	// Setup
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)

	// Test
	healthHandler := NewHealthHandler()
	healthHandler.Health(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
}
