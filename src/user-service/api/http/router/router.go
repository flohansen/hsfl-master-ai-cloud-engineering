package router

import (
	router "hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/middleware"
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
	userHandler user.Controller,
) *Router {
	r := router.New()

	r.POST("/api/v1/authentication/login/", loginHandler.Login)
	r.POST("/api/v1/authentication/register/", registerHandler.Register)

	authMiddleware := middleware.CreateAuthMiddleware()

	r.GET("/api/v1/user/", userHandler.GetUsers, authMiddleware)
	r.GET("/api/v1/user/role/:userRole", userHandler.GetUsersByRole, authMiddleware)
	r.GET("/api/v1/user/:userId", userHandler.GetUser, authMiddleware)
	r.PUT("/api/v1/user/:userId", userHandler.PutUser, authMiddleware)
	r.POST("/api/v1/user/", userHandler.PostUser, authMiddleware)
	r.DELETE("/api/v1/user/:userId", userHandler.DeleteUser, authMiddleware)

	return &Router{r}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
