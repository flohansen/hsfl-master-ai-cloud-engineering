package main

import (
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/http-router"
	"net/http"
)

func main() {
	router := http_router.New()
	router.GET("/users", func(w http.ResponseWriter, r *http.Request) {

	})

	router.GET("/users/:userId", func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId").(string)
	})

	router.PUT("/users/:userId", func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId").(string)
	})

	router.DELETE("/users/:userId", func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId").(string)
	})

	router.GET("/users/:userId/books", func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId").(string)
	})
}
