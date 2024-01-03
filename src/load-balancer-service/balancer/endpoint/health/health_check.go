package health

import (
	"log"
	"net"
	"net/url"
	"sync"
	"time"
)

type HealthCheck struct {
	url *url.URL

	mutex       sync.Mutex
	check       func(addr *url.URL) bool
	period      time.Duration
	cancel      chan struct{}
	isAvailable bool
}

type CheckFunction func(addr *url.URL) bool

func NewHealthCheck(url *url.URL, check CheckFunction, period time.Duration) *HealthCheck {
	hc := HealthCheck{url: url}
	hc.SetHealthCheck(check, period)
	return &hc
}

func (hc *HealthCheck) IsAvailable() bool {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()
	return hc.isAvailable
}

func (hc *HealthCheck) Stop() {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()
	hc.stopHealthCheck()
}

func (hc *HealthCheck) SetHealthCheck(check CheckFunction, period time.Duration) {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()

	hc.stopHealthCheck()
	hc.check = check
	hc.period = period
	hc.cancel = make(chan struct{})
	hc.isAvailable = hc.check(hc.url)
	hc.runHealthCheck()
}

func (hc *HealthCheck) runHealthCheck() {
	checkHealth := func() {
		hc.mutex.Lock()
		defer hc.mutex.Unlock()
		isAvailable := hc.check(hc.url)
		hc.isAvailable = isAvailable
	}

	go func() {
		t := time.NewTicker(hc.period)
		for {
			select {
			case <-t.C:
				checkHealth()
			case <-hc.cancel:
				t.Stop()
				return
			}
		}
	}()
}

func (hc *HealthCheck) stopHealthCheck() {
	if hc.cancel != nil {
		hc.cancel <- struct{}{}
		close(hc.cancel)
		hc.cancel = nil
	}
}

var DefaultHealthCheck = func(addr *url.URL) bool {
	conn, err := net.DialTimeout("tcp", addr.Host, 2*time.Second)
	if err != nil {
		log.Printf("Health check is missing heartbeat: " + addr.Host)
		return false
	}
	_ = conn.Close()
	return true
}
