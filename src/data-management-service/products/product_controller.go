package products

import (
	"errors"

	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/data-management-service/products/model"
)

type ProductsRepository struct {
	products map[uint64]*model.Product
}

func NewProductsRepository() *ProductsRepository {
	return &ProductsRepository{products: make(map[uint64]*model.Product)}
}

func (repo *ProductsRepository) Create(product *model.Product) (*model.Product, error) {
	_, found := repo.products[product.Id]
	if found {
		return nil, errors.New("product already exists")
	}
	repo.products[product.Id] = product

	return product, nil
}

func (repo *ProductsRepository) Delete(product *model.Product) error {
	_, found := repo.products[product.Id]
	if found {
		delete(repo.products, product.Id)
		return nil
	}

	return errors.New("product could not be deleted")
}

func (repo *ProductsRepository) FindById(id uint64) (*model.Product, error) {
	product, found := repo.products[id]
	if found {
		return product, nil
	}

	return nil, errors.New("user could not be found")
}

func (repo *ProductsRepository) Update(product *model.Product) (*model.Product, error) {
	product, foundError := repo.FindById(product.Id)

	if foundError != nil {
		return nil, errors.New("product can not be updated")
	}

	return product, nil
}
