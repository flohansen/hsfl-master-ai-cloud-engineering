package metrics

import (
	"fmt"
	"log"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
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

func (m *Metrics) DisplayMetrics() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	l := widgets.NewList()
	l.Title = "Metrics"
	l.Rows = []string{
		"Current Users: ...",
		"Total Requests: ...",
		"Successful Requests: ...",
		"Failed Requests: ...",
		"Average Response Time: ...",
		"Max Response Time: ...",
	}
	l.SetRect(0, 0, 50, 10)

	ui.Render(l)

	ticker := time.NewTicker(time.Millisecond * 50)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.Rows = []string{
				"Current Users: " + formatInt64(m.ActiveUsers),
				"Total Requests: " + formatInt64(m.TotalRequests),
				"Successful Requests: " + formatInt64(m.SuccessfulRequests),
				"Failed Requests: " + formatInt64(m.FailedRequests),
				"Average Response Time: " + m.GetAverageResponseTime().String(),
				"Max Response Time: " + m.MaxResponseTime.String(),
			}

			ui.Render(l)
		}
	}
}

func formatInt64(i int64) string {
	return fmt.Sprintf("%d", i)
}
