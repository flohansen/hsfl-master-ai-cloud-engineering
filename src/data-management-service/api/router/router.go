package router

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/data-management-service/prices"
	"net/http"

	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router"
)

type Router struct {
	router http.Handler
}

func New(priceController prices.Controller) *Router {
	r := router.New()

	r.GET("/api/v1/price/:productId/:userId", priceController.GetPrice)
	r.PUT("/api/v1/price/:productId/:userId", priceController.PutPrice)
	r.DELETE("/api/v1/price/:productId/:userId", priceController.DeletePrice)

	return &Router{r}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
