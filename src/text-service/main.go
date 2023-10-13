package main

import (
	"net/http"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
)

func main() {
	router := router.New()

	router.GET("/books", func(w http.ResponseWriter, r *http.Request) {

	})

	router.GET("/books/:bookId", func(w http.ResponseWriter, r *http.Request) {
		bookId := r.Context().Value("bookId").(string)
	})

	router.GET("/books/:bookId/chapters", func(w http.ResponseWriter, r *http.Request) {
		bookId := r.Context().Value("bookId").(string)
	})

	router.GET("/books/:bookId/chapters/:chapterId", func(w http.ResponseWriter, r *http.Request) {
		bookId := r.Context().Value("bookId").(string)
		chapterId := r.Context().Value("chapterId").(string)
	})

	router.POST("/books", func(w http.ResponseWriter, r *http.Request) {

	})

	router.PUT("/books/:bookId", func(w http.ResponseWriter, r *http.Request) {
		bookId := r.Context().Value("bookId").(string)
	})

	router.DELETE("/books/:bookId", func(w http.ResponseWriter, r *http.Request) {
		bookId := r.Context().Value("bookId").(string)
	})

	router.POST("/books/:bookId/chapters", func(w http.ResponseWriter, r *http.Request) {
		bookId := r.Context().Value("bookId").(string)
	})

	router.PUT("/books/:bookId/chapters/:chapterId", func(w http.ResponseWriter, r *http.Request) {
		bookId := r.Context().Value("bookId").(string)
		chapterId := r.Context().Value("chapterId").(string)
	})

	router.DELETE("/books/:bookId/chapters/:chapterId", func(w http.ResponseWriter, r *http.Request) {
		bookId := r.Context().Value("bookId").(string)
		chapterId := r.Context().Value("chapterId").(string)
	})
}
