package loadtest

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/benchmark/config"
	"gotest.tools/assert"
)

func TestParseURLs(t *testing.T) {
	t.Run("parse urls successfully", func(t *testing.T) {
		targets := []string{"http://localhost:8080", "http://localhost:8081", "http://localhost:8082"}

		urls, err := ParseURLs(targets)

		assert.NilError(t, err)
		assert.Equal(t, len(urls), 3)

		for i, u := range urls {
			assert.Equal(t, u.String(), targets[i])
		}
	})
}

func TestDoRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	url, _ := url.Parse(server.URL)

	t.Run("request is successful", func(t *testing.T) {
		code, err := DoRequest(url)

		assert.NilError(t, err)
		assert.Equal(t, code, uint64(200))
	})

	t.Run("request is not successful", func(t *testing.T) {
		server.Close()

		code, err := DoRequest(url)

		assert.Assert(t, err != nil)
		assert.Assert(t, code == uint64(0))
	})
}

func TestGetCurrentStage(t *testing.T) {
	specs := []config.Spec{
		{Duration: 8 * time.Second},
		{Duration: 12 * time.Second},
		{Duration: 20 * time.Second},
	}

	tests := []struct {
		name          string
		specs         []config.Spec
		elapsed       time.Duration
		expectedValue int
	}{
		{
			name:          "Test Case 1: No stages passed",
			specs:         specs,
			elapsed:       0,
			expectedValue: 0,
		},
		{
			name:          "Test Case 2: One stage passed",
			specs:         specs,
			elapsed:       9 * time.Second,
			expectedValue: 1,
		},
		{
			name:          "Test Case 3: Two stages passed",
			specs:         specs,
			elapsed:       23 * time.Second,
			expectedValue: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stage := GetCurrentStage(tt.specs, tt.elapsed)

			assert.Equal(t, stage, tt.expectedValue)
		})
	}
}

func TestCalculateOverallDuration(t *testing.T) {
	specs := []config.Spec{
		{Duration: 2 * time.Second},
		{Duration: 3 * time.Second},
		{Duration: 4 * time.Second},
	}

	duration := CalculateOverallDuration(specs)

	assert.Equal(t, duration, 9*time.Second)
}

func TestCalculateStageElapsedTime(t *testing.T) {
	specs := []config.Spec{
		{Duration: time.Duration(5) * time.Second},
		{Duration: time.Duration(10) * time.Second},
		{Duration: time.Duration(30) * time.Second},
	}

	tests := []struct {
		name          string
		specs         []config.Spec
		currentStage  int
		elapsed       time.Duration
		expectedValue time.Duration
	}{
		{
			name:          "Test Case 1: No stages passed",
			specs:         specs,
			currentStage:  0,
			elapsed:       time.Duration(30) * time.Second,
			expectedValue: time.Duration(30) * time.Second,
		},
		{
			name:          "Test Case 2: One stage passed",
			specs:         specs,
			currentStage:  1,
			elapsed:       time.Duration(7) * time.Second,
			expectedValue: time.Duration(2) * time.Second,
		},
		{
			name:          "Test Case 3: Two stages passed",
			specs:         specs,
			currentStage:  2,
			elapsed:       time.Duration(28) * time.Second,
			expectedValue: time.Duration(13) * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elapsed := CalculateStageElapsedTime(tt.specs, tt.currentStage, tt.elapsed)

			assert.Equal(t, elapsed, tt.expectedValue)
		})
	}
}

func TestCalculateSleepDuration(t *testing.T) {
	tests := []struct {
		name          string
		elapsed       time.Duration
		duration      time.Duration
		lastSleep     time.Duration
		targetSleep   time.Duration
		expectedValue time.Duration
	}{
		{
			name:          "Test Case 1: lastSleep < targetSleep",
			elapsed:       time.Duration(2) * time.Second,
			duration:      time.Duration(10) * time.Second,
			lastSleep:     time.Duration(1) * time.Second,
			targetSleep:   time.Duration(2) * time.Second,
			expectedValue: time.Duration(1.2 * float64(time.Second)),
		},
		{
			name:          "Test Case 2: lastSleep > targetSleep",
			elapsed:       time.Duration(2) * time.Second,
			duration:      time.Duration(10) * time.Second,
			lastSleep:     time.Duration(2) * time.Second,
			targetSleep:   time.Duration(1) * time.Second,
			expectedValue: time.Duration(1.8 * float64(time.Second)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sleepDuration := CalculateSleepDuration(tt.elapsed, tt.duration, tt.lastSleep, tt.targetSleep)

			assert.Equal(t, sleepDuration, tt.expectedValue)
		})
	}
}

func TestPickRandomURL(t *testing.T) {
	url1, _ := url.Parse("http://localhost:8080")
	url2, _ := url.Parse("http://localhost:8081")
	url3, _ := url.Parse("http://localhost:8082")

	targets := []*url.URL{url1, url2, url3}

	pickedUrl := PickRandomURL(targets)

	assert.Assert(t, pickedUrl == url1 || pickedUrl == url2 || pickedUrl == url3)
}
