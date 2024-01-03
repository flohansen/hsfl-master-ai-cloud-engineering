package products

import (
	"errors"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products/model"
	"sort"
)

type DemoRepository struct {
	products map[uint64]*model.Product
}

func NewDemoRepository() *DemoRepository {
	return &DemoRepository{products: make(map[uint64]*model.Product)}
}

func (repo *DemoRepository) Create(product *model.Product) (*model.Product, error) {
	if product.Id == 0 {
		product.Id = repo.findNextAvailableID()
	}

	if foundProductWithEan, _ := repo.FindByEan(product.Ean); foundProductWithEan != nil {
		return nil, errors.New(ErrorEanAlreadyExists)
	}

	if _, found := repo.products[product.Id]; found {
		return nil, errors.New(ErrorProductAlreadyExists)
	}

	repo.products[product.Id] = product

	return product, nil
}

func (repo *DemoRepository) Delete(product *model.Product) error {
	_, found := repo.products[product.Id]
	if found {
		delete(repo.products, product.Id)
		return nil
	}

	return errors.New(ErrorProductDeletion)
}

func (repo *DemoRepository) FindAll() ([]*model.Product, error) {
	if repo.products != nil {
		r := make([]*model.Product, 0, len(repo.products))
		for _, v := range repo.products {
			r = append(r, v)
		}

		sort.Slice(r, func(i, j int) bool {
			return r[i].Description < r[j].Description
		})
		return r, nil
	}

	return nil, errors.New(ErrorProductsList)
}

func (repo *DemoRepository) FindById(id uint64) (*model.Product, error) {
	product, found := repo.products[id]
	if found {
		return product, nil
	}

	return nil, errors.New(ErrorProductNotFound)
}

func (repo *DemoRepository) FindByEan(ean uint64) (*model.Product, error) {
	for _, entry := range repo.products {
		if entry.Ean == ean {
			return entry, nil
		}
	}

	return nil, errors.New(ErrorProductNotFound)
}

func (repo *DemoRepository) Update(product *model.Product) (*model.Product, error) {
	existingProduct, foundError := repo.FindById(product.Id)

	if foundError != nil {
		return nil, errors.New(ErrorProductUpdate)
	}

	existingProduct.Description = product.Description
	existingProduct.Ean = product.Ean

	return existingProduct, nil
}

func (repo *DemoRepository) findNextAvailableID() uint64 {
	var maxID uint64
	for id := range repo.products {
		if id > maxID {
			maxID = id
		}
	}
	return maxID + 1
}
