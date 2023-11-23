package products

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products/model"
)

type Repository interface {
	Create(*model.Product) (*model.Product, error)
	Delete(*model.Product) error
	FindAll() ([]*model.Product, error)
	FindById(id uint64) (*model.Product, error)
	FindByEan(id uint64) ([]*model.Product, error)
	Update(*model.Product) (*model.Product, error)
}

const (
	ErrorProductsList         = "product list not available"
	ErrorProductNotFound      = "product could not be found"
	ErrorProductUpdate        = "product can not be updated"
	ErrorProductDeletion      = "product could not be deleted"
	ErrorProductAlreadyExists = "product already exists"
)
