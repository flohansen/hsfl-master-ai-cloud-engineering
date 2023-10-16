package main

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/data-management-service/api/router"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/data-management-service/products"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/data-management-service/products/model"
	"log"
	"net/http"
)

func main() {
	productRepository := products.NewDemoRepository()
	productsController := products.NewDefaultController(productRepository)

	productSlice := []*model.Product{
		{
			Id:          1,
			Description: "Strauchtomaten",
			Ean:         4014819040771,
		},
		{
			Id:          2,
			Description: "Lauchzwiebeln",
			Ean:         5001819040871,
		},
	}

	for _, product := range productSlice {
		_, err := productRepository.Create(product)
		if err != nil {
			return
		}
	}

	handler := router.New(productsController)

	if err := http.ListenAndServe(":3001", handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
