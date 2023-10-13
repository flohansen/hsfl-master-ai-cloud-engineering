package main

import (
	"net/http"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
)

func main() {
	router := router.New()

	// Authentication stuff
	router.POST("/login", func(w http.ResponseWriter, r *http.Request) {

	})

	router.POST("/register", func(w http.ResponseWriter, r *http.Request) {

	})

	router.GET("/refresh-token", func(w http.ResponseWriter, r *http.Request) {

	})

	router.GET("/logout", func(w http.ResponseWriter, r *http.Request) {

	})

	// user-specific stuff
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
