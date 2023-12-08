package loadtest

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestCalculateSleepDuration(t *testing.T) {
	rampupDuration := time.Duration(10) * time.Second

	t.Run("elapsed time is less than rampupDuration", func(t *testing.T) {
		elapsed := time.Duration(5) * time.Second
		expectedSleepDuration := MinimumSleepDuration + time.Duration((1.0-(float64(elapsed)/float64(rampupDuration)))*DecreaseInterval)
		actualSleepDuration := CalculateSleepDuration(elapsed, rampupDuration)

		assert.Equal(t, expectedSleepDuration, actualSleepDuration)
	})

	t.Run("elapsed time is more than rampupDuration", func(t *testing.T) {
		elapsed := time.Duration(15) * time.Second
		expectedSleepDuration := MinimumSleepDuration
		actualSleepDuration := CalculateSleepDuration(elapsed, rampupDuration)

		assert.Equal(t, expectedSleepDuration, actualSleepDuration)
	})
}

func TestPickRandomTarget(t *testing.T) {
	targets := []string{"target1", "target2", "target3"}

	t.Run("should return a random target", func(t *testing.T) {
		target := PickRandomTarget(targets)

		assert.Assert(t, target == "target1" || target == "target2" || target == "target3")
	})
}

func TestUpdateResponseTimes(t *testing.T) {
	target := "target1"
	requestStart := time.Now()
	responseTimesByTarget := make(map[string]responseTimeEntry)
	mu := &sync.Mutex{}

	UpdateResponseTimes(target, requestStart, responseTimesByTarget, mu)

	assert.Equal(t, 1, len(responseTimesByTarget))

	data, ok := responseTimesByTarget[target]

	assert.Assert(t, ok)
	assert.Equal(t, 1, data.count)
	assert.Assert(t, data.total > 0)
}

func TestMakeRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	target := server.URL
	id := 1
	responseTimesByTarget := make(map[string]responseTimeEntry)
	mu := &sync.Mutex{}

	err := MakeRequest(target, id, responseTimesByTarget, mu)

	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(responseTimesByTarget))

	data, ok := responseTimesByTarget[target]

	assert.Assert(t, ok)
	assert.Equal(t, 1, data.count)
	assert.Assert(t, data.total > 0)
}

func TestPrintAverageResponseTimes(t *testing.T) {
	responseTimesByTarget := make(map[string]responseTimeEntry)
	target := "http://example.com"
	responseTimesByTarget[target] = responseTimeEntry{
		total: time.Duration(10) * time.Second,
		count: 5,
	}

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(nil)
	}()

	PrintAverageResponseTimes(responseTimesByTarget)

	expected := "Average response time for  http://example.com  was  2s"
	if !strings.Contains(buf.String(), expected) {
		t.Errorf("Expected log output to contain %v, but got %v", expected, buf.String())
	}
}
