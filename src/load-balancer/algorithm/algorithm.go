package algorithm

import (
	"net/http"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/load-balancer/model"
)

type Algorithm interface {
	GetTarget(_ *http.Request, targets []model.Target, fun func(target model.Target))
}
