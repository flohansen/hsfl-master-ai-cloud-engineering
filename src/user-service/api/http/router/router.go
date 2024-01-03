package router

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/api/http/handler"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user"
	"net/http"
)

type Router struct {
	router http.Handler
}

func New(
	loginHandler *handler.LoginHandler,
	registerHandler *handler.RegisterHandler,
	userController *user.Controller,
	authMiddleware router.Middleware,
) *Router {
	r := router.New()

	r.POST("/api/v1/authentication/login/", loginHandler.Login)
	r.POST("/api/v1/authentication/register/", registerHandler.Register)

	r.GET("/api/v1/user/role/:userRole", (*userController).GetUsersByRole, authMiddleware)
	r.GET("/api/v1/user/:userId", (*userController).GetUser, authMiddleware)
	r.PUT("/api/v1/user/:userId", (*userController).PutUser, authMiddleware)
	r.DELETE("/api/v1/user/:userId", (*userController).DeleteUser, authMiddleware)

	return &Router{r}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
