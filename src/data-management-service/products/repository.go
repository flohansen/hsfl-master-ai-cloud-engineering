package products

import "hsfl.de/group6/hsfl-master-ai-cloud-engineering/data-management-service/products/model"

type Repository interface {
	Create(*model.Product) (*model.Product, error)
	Delete(*model.Product) error
	FindAll() ([]*model.Product, error)
	FindById(id uint64) (*model.Product, error)
	Update(*model.Product) (*model.Product, error)
}
