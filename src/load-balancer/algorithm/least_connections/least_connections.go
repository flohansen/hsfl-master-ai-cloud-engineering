package least_connections

import (
	"math"
	"net/http"
	"sync"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/load-balancer/model"
)

type LeastConnections struct {
	mutex sync.RWMutex
	count map[model.Target]int
}

func New() *LeastConnections {
	return &LeastConnections{sync.RWMutex{}, make(map[model.Target]int)}
}

func (algorithm *LeastConnections) GetTarget(r *http.Request, targets []model.Target, fun func(target model.Target)) {
	algorithm.mutex.RLock()

	var targetHost model.Target

	minConnections := math.MaxInt

	for _, host := range targets {
		if currentConnections, ok := algorithm.count[host]; ok && currentConnections < minConnections {
			minConnections = currentConnections
			targetHost = host
		}
	}

	algorithm.mutex.RUnlock()

	algorithm.adjustConnectionCount(targetHost, 1)

	fun(targetHost)

	algorithm.adjustConnectionCount(targetHost, -1)
}

func (algorithm *LeastConnections) adjustConnectionCount(host model.Target, delta int) {
	algorithm.mutex.Lock()
	defer algorithm.mutex.Unlock()

	algorithm.count[host] += delta
}
