package main

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/web-service/api/router"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/web-service/service"
	"log"
	"net/http"
)

func main() {
	productsController := service.NewDefaultController()
	handler := router.New(productsController)

	if err := http.ListenAndServe(":3000", handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
