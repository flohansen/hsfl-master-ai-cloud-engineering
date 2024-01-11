package prices

import "hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices/model"

type Repository interface {
	Create(*model.Price) (*model.Price, error)
	FindAll() ([]*model.Price, error)
	FindAllByUser(userId uint64) ([]*model.Price, error)
	FindAllByProduct(productId uint64) ([]*model.Price, error)
	FindByIds(productId uint64, userId uint64) (*model.Price, error)
	Update(*model.Price) (*model.Price, error)
	Delete(*model.Price) error
}

const (
	ErrorPriceList          = "price list not available"
	ErrorPriceNotFound      = "price could not be found"
	ErrorPriceUpdate        = "price can not be updated"
	ErrorPriceDeletion      = "price could not be deleted"
	ErrorPriceAlreadyExists = "price already exists"
)
