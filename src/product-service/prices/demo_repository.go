package prices

import (
	"errors"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices/model"
	"reflect"
)

type priceEntryKey struct {
	UserId    uint64
	ProductId uint64
}

type DemoRepository struct {
	prices map[priceEntryKey]*model.Price
}

func NewDemoRepository() *DemoRepository {
	return &DemoRepository{prices: make(map[priceEntryKey]*model.Price)}
}

func (repo *DemoRepository) Create(price *model.Price) (*model.Price, error) {
	for _, p := range repo.prices {
		if reflect.DeepEqual(p, price) {
			return nil, errors.New(ErrorPriceAlreadyExists)
		}
	}

	key := priceEntryKey{
		ProductId: price.ProductId,
		UserId:    price.UserId,
	}

	repo.prices[key] = price

	return price, nil
}

func (repo *DemoRepository) Delete(price *model.Price) error {
	key := priceEntryKey{
		UserId:    price.UserId,
		ProductId: price.ProductId,
	}

	_, exists := repo.prices[key]
	if !exists {
		return errors.New(ErrorPriceDeletion)
	}

	delete(repo.prices, key)
	return nil
}

func (repo *DemoRepository) FindAll() ([]*model.Price, error) {
	if repo.prices != nil {
		r := make([]*model.Price, 0, len(repo.prices))
		for _, v := range repo.prices {
			r = append(r, v)
		}
		return r, nil
	}

	return nil, errors.New(ErrorPriceList)
}

func (repo *DemoRepository) FindByIds(productId uint64, userId uint64) (*model.Price, error) {
	key := priceEntryKey{
		UserId:    userId,
		ProductId: productId,
	}

	price, exists := repo.prices[key]

	if !exists {
		return nil, errors.New(ErrorPriceNotFound)
	}

	return price, nil
}

func (repo *DemoRepository) Update(price *model.Price) (*model.Price, error) {
	existingPrice, foundError := repo.FindByIds(price.ProductId, price.UserId)

	if foundError != nil {
		return nil, errors.New(ErrorPriceUpdate)
	}

	existingPrice.Price = price.Price

	return existingPrice, nil
}
