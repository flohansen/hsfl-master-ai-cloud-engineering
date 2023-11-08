package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"
)

func TestReverseProxyHandler_ServeHTTP(t *testing.T) {
	config := ReverseProxyConfig{
		AuthServiceURL:         "http://authserviceurl",
		BulletinBoardServiceURL: "http://bulletinboardserviceurl",
		FeedServiceURL:         "http://localhost:8081",
	}
	reverseProxy := NewReverseProxy(config)

	// Create a test request
	req, err := http.NewRequest("GET", "http://localhost:8080/feed/feed", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test response recorder
	recorder := httptest.NewRecorder()

	// Call the ServeHTTP method to process the request
	reverseProxy.ServeHTTP(recorder, req)

	// Check the response status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	// You can add more test cases for other scenarios as needed
}

func TestReverseProxyHandler_determineTarget(t *testing.T) {
	reverseProxy := &ReverseProxyHandler{}

	// Test a request with a valid target
	path := "/feed/feed"
	target := reverseProxy.determineTarget(path)
	if target != "feed" {
		t.Errorf("Expected target 'feed', got '%s'", target)
	}

	// Test a request with an invalid target
	path = "/unknown/somepath"
	target = reverseProxy.determineTarget(path)

	if target != "unknown" {
		t.Errorf("Expected target 'unknown', got '%s'", target)
	}

	// Test a request with an empty path
	path = ""
	target = reverseProxy.determineTarget(path)
	fmt.Println(target)

	if target != "unknown" {
		t.Errorf("Expected target 'unknown', got '%s'", target)
	}

	// You can add more test cases for other scenarios as needed
}