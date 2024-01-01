package tester

import (
	"http-stress-test/config"
	"http-stress-test/metrics"
	"http-stress-test/network"
	"math/rand"
	"sync"
	"time"
)

type tester struct {
	config  *config.Configuration
	metrics *metrics.Metrics
}

func NewTester(config *config.Configuration, metrics *metrics.Metrics) *tester {
	return &tester{
		config:  config,
		metrics: metrics,
	}
}

func (t *tester) Run() {
	var wg sync.WaitGroup

	rampUpInterval := time.Duration((int64(t.config.RampUp) * time.Second.Nanoseconds()) / int64(t.config.Users))

	for i := 0; i < t.config.Users; i++ {
		wg.Add(1)
		go t.runUser(&wg, t.metrics)

		if t.config.RampUp > 0 {
			time.Sleep(rampUpInterval)
		}
	}

	wg.Wait()
	time.Sleep(time.Second)
}

func (t *tester) runUser(wg *sync.WaitGroup, metrics *metrics.Metrics) {
	if t.metrics != nil {
		// Metrics
		metrics.IncrementUserCount()
		defer metrics.DecrementUserCount()
	}
	// Send done to waiting group
	defer wg.Done()

	client := network.NewHttpClient()

	// Calculate force ending time
	var endTime time.Time
	if t.config.Duration > 0 {
		endTime = time.Now().Add(time.Duration(t.config.Duration) * time.Second)
	}

	requestCount := 0
	for {
		// Cancel on max duration
		if !endTime.IsZero() && time.Now().After(endTime) {
			break
		}

		// Cancel on max requests
		if t.config.Requests > 0 && requestCount >= t.config.Requests {
			break
		}

		// Random target
		targetIndex := rand.Intn(len(t.config.Targets))
		targetURL := t.config.Targets[targetIndex].URL

		// Send and collect metric
		startTime := time.Now()
		response, err := client.SendRequest(targetURL)
		responseTime := time.Since(startTime)

		// Record metric
		if metrics != nil {
			success := err == nil && response.StatusCode() == 200
			metrics.RecordResponse(responseTime, success)
		}

		// Limiting rate if configured
		if t.config.RateLimit > 0 {
			time.Sleep(time.Second / time.Duration(t.config.RateLimit))
		}

		requestCount++
	}
}
