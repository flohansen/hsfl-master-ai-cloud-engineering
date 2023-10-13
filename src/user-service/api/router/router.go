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

	// CRUD User

	// Current User Info
	r.GET("/api/v1/user/", currentUserInfoHandler)
	// User Info
	r.GET("/api/v1/user/:userId", userInfoHandler)
	// Register User
	r.PUT("/api/v1/user/register", registerHandler)
	// Authenticate User
	r.POST("/api/v1/user/auth", authHandler)
	// Delete User
	r.DELETE("/api/v1/user/delete", authHandler)

	return &Router{r}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
