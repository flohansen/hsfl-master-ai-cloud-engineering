package router

import (
	router "hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/web-service/service"
	"net/http"
)

type Router struct {
	router http.Handler
}

func New(
	serviceController service.Controller,
) *Router {
	r := router.New()

	r.GET("/user/shoppingList", serviceController.GetShoppingList)
	r.GET("/admin/products", serviceController.GetAdminProducts)
	r.GET("/merchant/products", serviceController.GetMerchantProducts)
	r.GET("/productCatalogue", serviceController.GetProductCatalogue)

	return &Router{r}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
