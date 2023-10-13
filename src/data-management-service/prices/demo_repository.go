package prices

import (
	"errors"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/data-management-service/prices/model"
	"reflect"
)

type DemoRepository struct {
	prices []*model.Price
}

func NewDemoRepository() *DemoRepository {
	return &DemoRepository{prices: make([]*model.Price, 1)}
}

func (repo *DemoRepository) Create(price *model.Price) (*model.Price, error) {
	repo.prices = append(repo.prices, price)
	return price, nil
}

func (repo *DemoRepository) Delete(priceToDelete *model.Price) error {
	for i, price := range repo.prices {
		if reflect.DeepEqual(price, priceToDelete) {
			repo.prices = append(repo.prices[:i], repo.prices[i+1:]...)
			return nil
		}
	}

	return errors.New("price could not be deleted")
}

func (repo *DemoRepository) FindByIds(productId uint64, userId uint64) (*model.Price, error) {
	for _, price := range repo.prices {
		if price.ProductId == productId && price.UserId == userId {
			return price, nil
		}
	}

	return nil, errors.New("price could not be found")
}

func (repo *DemoRepository) Update(price *model.Price) (*model.Price, error) {
	price, foundError := repo.FindByIds(price.ProductId, price.UserId)

	if foundError != nil {
		return nil, errors.New("price can not be updated")
	}

	return price, nil
}
