package router

import (
	router "hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/api/handler"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user"
	"net/http"
)

type Router struct {
	router http.Handler
}

func New(
	loginHandler *handler.LoginHandler,
	registerHandler *handler.RegisterHandler,
	userHandler user.Controller,
) *Router {
	r := router.New()

	r.POST("/api/v1/user/login", loginHandler.Login)
	r.POST("/api/v1/user/register", registerHandler.Register)

	r.GET("/api/v1/user/:userId", userHandler.GetUser)

	return &Router{r}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
