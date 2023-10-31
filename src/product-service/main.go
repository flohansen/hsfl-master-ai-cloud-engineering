package main

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/api/router"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices"
	priceModel "hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices/model"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products"
	productModel "hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products/model"
	"log"
	"net/http"
)

func main() {
	productRepository := products.NewDemoRepository()
	productsController := products.NewDefaultController(productRepository)
	createContentForProducts(productRepository)

	priceRepository := prices.NewDemoRepository()
	pricesController := prices.NewDefaultController(priceRepository)
	createContentForPrices(priceRepository)

	handler := router.New(productsController, pricesController)

	if err := http.ListenAndServe(":3001", handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}

func createContentForPrices(priceRepository prices.Repository) {
	pricesSlice := []*priceModel.Price{
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

	for _, price := range pricesSlice {
		_, err := priceRepository.Create(price)
		if err != nil {
			return
		}
	}
}

func createContentForProducts(productRepository products.Repository) {
	productSlice := []*productModel.Product{
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
}
