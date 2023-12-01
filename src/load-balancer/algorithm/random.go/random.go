package random

import (
	"math/rand"
	"net/http"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/load-balancer/model"
)

type Random struct{}

func New() *Random {
	return &Random{}
}

func (algorithm *Random) GetTarget(r *http.Request, targets []model.Target, fun func(target model.Target)) {
	if len(targets) == 0 {
		return
	}

	target := targets[rand.Intn(len(targets))]

	fun(target)
}
