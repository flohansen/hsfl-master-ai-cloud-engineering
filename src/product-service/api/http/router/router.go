package router

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products"
	"net/http"
)

type Router struct {
	router http.Handler
}

func New(
	productController *products.Controller,
	pricesController *prices.Controller,
	authMiddleware router.Middleware,
) *Router {
	r := router.New()

	r.GET("/api/v1/product/", (*productController).GetProducts)
	r.GET("/api/v1/product/ean/:productEan", (*productController).GetProductByEan)
	r.GET("/api/v1/product/:productId", (*productController).GetProductById)
	r.PUT("/api/v1/product/:productId", (*productController).PutProduct, authMiddleware)
	r.POST("/api/v1/product/", (*productController).PostProduct, authMiddleware)
	r.DELETE("/api/v1/product/:productId", (*productController).DeleteProduct, authMiddleware)

	r.GET("/api/v1/price/", (*pricesController).GetPrices)
	r.GET("/api/v1/price/user/:userId", (*pricesController).GetPricesByUser)
	r.GET("/api/v1/price/product/:productId", (*pricesController).GetPricesByProduct)
	r.GET("/api/v1/price/:productId/:userId", (*pricesController).GetPrice)
	r.PUT("/api/v1/price/:productId/:userId", (*pricesController).PutPrice, authMiddleware)
	r.POST("/api/v1/price/:productId/:userId", (*pricesController).PostPrice, authMiddleware)
	r.DELETE("/api/v1/price/:productId/:userId", (*pricesController).DeletePrice, authMiddleware)

	return &Router{r}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
