package router

import (
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/user"
	"net/http"
)

type Router struct {
	router http.Handler
}

func New(
	userController user.Controller,
) *Router {
	router := router.New()

	router.POST("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		userController.Login(w, r)
	})
	router.POST("/api/v1/register", func(w http.ResponseWriter, r *http.Request) {
		userController.Register(w, r)
	})

	router.GET("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		userController.GetUsers(w, r)
	})
	router.GET("/api/v1/users/:username", func(w http.ResponseWriter, r *http.Request) {
		userController.GetUser(w, r)
	})
	router.PUT("/api/v1/users/:username", func(w http.ResponseWriter, r *http.Request) {
		userController.PutUser(w, r)
	})
	router.DELETE("/api/v1/users/:username", func(w http.ResponseWriter, r *http.Request) {
		userController.DeleteUser(w, r)
	})

	return &Router{router}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
