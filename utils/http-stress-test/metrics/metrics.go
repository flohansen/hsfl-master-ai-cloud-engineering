package metrics

import (
	"context"
	"fmt"
	"github.com/pterm/pterm"
	"sync"
	"time"
)

type Metrics struct {
	ActiveUsers        int64
	TotalRequests      int64
	SuccessfulRequests int64
	FailedRequests     int64
	TotalResponseTime  time.Duration
	MaxResponseTime    time.Duration
	lock               sync.Mutex
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) IncrementUserCount() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.ActiveUsers++
}

func (m *Metrics) DecrementUserCount() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.ActiveUsers--
}

func (m *Metrics) RecordResponse(duration time.Duration, success bool) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.TotalRequests++
	if success {
		m.SuccessfulRequests++
	} else {
		m.FailedRequests++
	}
	m.TotalResponseTime += duration
	if duration > m.MaxResponseTime {
		m.MaxResponseTime = duration
	}
}

func (m *Metrics) GetAverageResponseTime() time.Duration {
	if m.TotalRequests == 0 {
		return 0
	}
	return time.Duration(m.TotalResponseTime.Nanoseconds() / m.TotalRequests)
}

func (m *Metrics) GetAverageRequestsPerSecond(startTime time.Time) float64 {
	elapsedTime := time.Since(startTime).Seconds()
	if elapsedTime == 0 {
		return 0
	}
	return float64(m.TotalRequests) / elapsedTime
}

func (m *Metrics) DisplayMetrics(ctx context.Context) {
	startTime := time.Now()

	// Initialize pterm
	pterm.EnableColor()

	area, _ := pterm.DefaultArea.WithFullscreen(true).Start()
	defer area.Stop()

	updateMetricsData := func() [][]string {
		return [][]string{
			{"HTTP Stress Test", ""},
			{"Current Users", fmt.Sprintf("%d", m.ActiveUsers)},
			{"Total Requests", fmt.Sprintf("%d", m.TotalRequests)},
			{"Successful Requests", fmt.Sprintf("%d", m.SuccessfulRequests)},
			{"Failed Requests", fmt.Sprintf("%d", m.FailedRequests)},
			{"Average Response Time", m.GetAverageResponseTime().String()},
			{"Max Response Time", m.MaxResponseTime.String()},
			{"Avg Requests/Sec", fmt.Sprintf("%.2f", m.GetAverageRequestsPerSecond(startTime))},
		}
	}

	// Ticker for updating the UI
	ticker := time.NewTicker(time.Millisecond * 50)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Create and render the table
			data := updateMetricsData()
			table, _ := pterm.DefaultTable.WithHasHeader().WithData(data).Srender()
			area.Update(table)
		case <-ctx.Done():
			return
		}
	}
}

func formatInt64(i int64) string {
	return fmt.Sprintf("%d", i)
}
