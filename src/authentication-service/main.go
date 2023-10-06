package main

import (
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/http-router"
	"net/http"
)

func main() {
	router := http_router.New()
	router.POST("/login", func(w http.ResponseWriter, r *http.Request) {

	})

	router.POST("/register", func(w http.ResponseWriter, r *http.Request) {

	})

	router.GET("/refresh-token", func(w http.ResponseWriter, r *http.Request) {

	})

	router.GET("/logout", func(w http.ResponseWriter, r *http.Request) {

	})
}
