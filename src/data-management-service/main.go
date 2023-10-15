package main

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/data-management-service/api/router"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/data-management-service/prices"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/data-management-service/prices/model"
	"log"
	"net/http"
)

func main() {
	priceRepository := prices.NewDemoRepository()
	pricesController := prices.NewDefaultController(priceRepository)

	pricesSlice := []*model.Price{
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

	handler := router.New(pricesController)

	if err := http.ListenAndServe(":3001", handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
