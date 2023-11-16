package round_robin

import (
	"net/http"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/load-balancer/model"
)

type RoundRobin struct {
	idx int
}

func New() *RoundRobin {
	return &RoundRobin{0}
}

func (algorithm *RoundRobin) GetTarget(_ *http.Request, targets []model.Target, fun func(replica model.Target)) {
	if len(targets) == 0 {
		return
	}

	target := targets[algorithm.idx]
	algorithm.idx = (algorithm.idx + 1) % len(targets)
	fun(target)
}
