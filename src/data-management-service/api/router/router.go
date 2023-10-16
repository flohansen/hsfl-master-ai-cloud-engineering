package router

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/data-management-service/products"
	"net/http"

	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router"
)

type Router struct {
	router http.Handler
}

func New(productController products.Controller) *Router {
	r := router.New()

	r.GET("/api/v1/product/", productController.GetProducts)
	r.GET("/api/v1/product/:productId", productController.GetProduct)

	r.PUT("/api/v1/product/:productId", productController.PutProduct)

	r.DELETE("/api/v1/product/:productId", productController.DeleteProduct)
	return &Router{r}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
