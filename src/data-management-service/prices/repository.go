package prices

import "hsfl.de/group6/hsfl-master-ai-cloud-engineering/data-management-service/prices/model"

type Repository interface {
	Create(*model.Price) (*model.Price, error)
	Delete(*model.Price) error
	FindByIds(productId uint64, userId uint64) (*model.Price, error)
	Update(*model.Price) (*model.Price, error)
}
