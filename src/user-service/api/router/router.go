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

	router.POST("/api/v1/login", userController.Login)
	router.POST("/api/v1/register", userController.Register)

	router.GET("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println("Hallo, ich bin hiiiiiiiiiiiiiiiiiiiier", r, r.Method)
		userController.GetUsers(w, r)
	})
	router.GET("/api/v1/users/:username", userController.GetUser)
	router.PUT("/api/v1/users/:username", userController.PutUser)
	router.DELETE("/api/v1/users/:username", userController.DeleteUser)

	return &Router{router}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
