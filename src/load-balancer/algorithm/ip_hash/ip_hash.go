package ip_hash

import (
	"hash/fnv"
	"net/http"
	"regexp"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/load-balancer/model"
)

type IpHash struct{}

func New() *IpHash {
	return &IpHash{}
}

func (algorithm *IpHash) GetTarget(r *http.Request, targets []model.Target, fun func(target model.Target)) {
	if len(targets) == 0 {
		return
	}

	ip := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`).FindString(r.RemoteAddr)

	hashValue := hash(ip)
	target := targets[hashValue%uint32(len(targets))]

	fun(target)
}

func hash(str string) uint32 {
	hasher := fnv.New32a()
	hasher.Write([]byte(str))

	return hasher.Sum32()
}
