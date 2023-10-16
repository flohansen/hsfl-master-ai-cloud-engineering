package router

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/data-management-service/prices"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/data-management-service/products"
	"net/http"

	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router"
)

type Router struct {
	router http.Handler
}

func New(productController products.Controller, pricesController prices.Controller) *Router {
	r := router.New()

	r.GET("/api/v1/product/", productController.GetProducts)
	r.GET("/api/v1/product/:productId", productController.GetProduct)
	r.PUT("/api/v1/product/:productId", productController.PutProduct)
	r.DELETE("/api/v1/product/:productId", productController.DeleteProduct)

	r.GET("/api/v1/price/:productId/:userId", pricesController.GetPrice)
	r.PUT("/api/v1/price/:productId/:userId", pricesController.PutPrice)
	r.DELETE("/api/v1/price/:productId/:userId", pricesController.DeletePrice)

	return &Router{r}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
