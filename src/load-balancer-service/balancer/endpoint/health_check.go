package endpoint

import (
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

func (hc *HealthCheck) SetHealthCheck(check func(addr *url.URL) bool, period time.Duration) {
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

var defaultHealthCheck = func(addr *url.URL) bool {
	conn, err := net.DialTimeout("tcp", addr.Host, 10)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
