package router

import (
	"net/http"

	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router"
)

type Router struct {
	router http.Handler
}

func New() *Router {
	r := router.New()

	// r.GET("/api/v1/products")
	// r.POST("/api/v1/products")
	// r.GET("/api/v1/products/:productid")
	// r.PUT("/api/v1/products/:productid")
	// r.DELETE("/api/v1/products/:productid")

	return &Router{r}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
