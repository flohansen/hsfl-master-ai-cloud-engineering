package router

import (
	"net/http"

	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router"
)

type Router struct {
	router http.Handler
}

func New() *Router {
	router := router.New()

	// router.GET("/api/v1/products")
	// router.POST("/api/v1/products")
	// router.GET("/api/v1/products/:productid")
	// router.PUT("/api/v1/products/:productid")
	// router.DELETE("/api/v1/products/:productid")

	return &Router{router}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
