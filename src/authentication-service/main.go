package main

import (
	"net/http"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
)

func main() {
	router := router.New()
	router.POST("/login", func(w http.ResponseWriter, r *http.Request) {

	})

	router.POST("/register", func(w http.ResponseWriter, r *http.Request) {

	})

	router.GET("/refresh-token", func(w http.ResponseWriter, r *http.Request) {

	})

	router.GET("/logout", func(w http.ResponseWriter, r *http.Request) {

	})
}
