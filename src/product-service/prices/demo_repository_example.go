package prices

import "hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices/model"

func GenerateExampleDemoRepository() Repository {
	repository := NewDemoRepository()
	pricesSlice := GenerateExamplePriceSlice()
	for _, price := range pricesSlice {
		repository.Create(price)
	}

	return repository
}

func GenerateExamplePriceSlice() []*model.Price {
	return []*model.Price{
		{
			UserId:    1,
			ProductId: 1,
			Price:     2.99,
		},
		{
			UserId:    2,
			ProductId: 2,
			Price:     5.99,
		},
	}
}
