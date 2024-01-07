package products

import "hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products/model"

func GenerateExampleDemoRepository() Repository {
	repository := NewDemoRepository()
	productSlice := GenerateExampleProductSlice()
	for _, product := range productSlice {
		repository.Create(product)
	}

	return repository
}

func GenerateExampleProductSlice() []*model.Product {
	return []*model.Product{
		{
			Id:          1,
			Description: "Strauchtomaten",
			Ean:         "4014819040771",
		},
		{
			Id:          2,
			Description: "Lauchzwiebeln",
			Ean:         "5001819040871",
		},
	}
}
