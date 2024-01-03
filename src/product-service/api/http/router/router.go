package router

import (
	router "hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products"
	"net/http"
)

type Router struct {
	router http.Handler
}

func New(productController *products.Controller, pricesController *prices.Controller) *Router {
	r := router.New()

	r.GET("/api/v1/product/", (*productController).GetProducts)
	r.GET("/api/v1/product/ean/:productEan", (*productController).GetProductByEan)
	r.GET("/api/v1/product/:productId", (*productController).GetProductById)
	r.PUT("/api/v1/product/:productId", (*productController).PutProduct)
	r.POST("/api/v1/product/", (*productController).PostProduct)
	r.DELETE("/api/v1/product/:productId", (*productController).DeleteProduct)

	r.GET("/api/v1/price/", (*pricesController).GetPrices)
	r.GET("/api/v1/price/user/:userId", (*pricesController).GetPricesByUser)
	r.GET("/api/v1/price/:productId/:userId", (*pricesController).GetPrice)
	r.PUT("/api/v1/price/:productId/:userId", (*pricesController).PutPrice)
	r.POST("/api/v1/price/:productId/:userId", (*pricesController).PostPrice)
	r.DELETE("/api/v1/price/:productId/:userId", (*pricesController).DeletePrice)

	return &Router{r}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
