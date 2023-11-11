package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServeHTTP(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send response to be tested
		rw.Write([]byte(`OK`))
	}))

	// Close the server when the test finishes
	defer server.Close()

	service := Service{
		Name:        "test-service",
		ContextPath: "/service",
		TargetURL:   server.URL,
	}

	config := ReverseProxyConfig{Services: []Service{service}}
	rp := NewReverseProxy(config)

	req := httptest.NewRequest("GET", "/service/path", nil)
	w := httptest.NewRecorder()

	rp.ServeHTTP(w, req)
	res := w.Result()

	// Check the status code is what we expect
	assert.Equal(t, http.StatusOK, res.StatusCode, "Expected response status 200 OK")

	// Check the response body is what we expect
	expected := `OK`
	actual := w.Body.String()
	assert.Equal(t, expected, actual)
}

func TestServeHTTP_NotFound(t *testing.T) {
	rp := &ReverseProxyHandler{}

	req := httptest.NewRequest("GET", "/nonexistent/path", nil)
	w := httptest.NewRecorder()

	rp.ServeHTTP(w, req)

	res := w.Result()
	// Check the status code is what we expect
	assert.Equal(t, http.StatusNotFound, res.StatusCode, "Expected response status 404 Not Found")
}
