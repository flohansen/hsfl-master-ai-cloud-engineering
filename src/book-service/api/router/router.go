package router

import (
	"fmt"
	"net/http"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
)

type Router struct {
	router http.Handler
}

func New(
	booksController books.Controller,
) *Router {
	router := router.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request) {

	})

	router.GET("/api/v1/books", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("router")
		fmt.Println(w, r)
		booksController.GetBooks(w, r)
	})
	router.POST("/api/v1/books", booksController.PostBooks)
	router.GET("/api/v1/books/:bookid", booksController.GetBook)
	router.PUT("/api/v1/books/:bookid", booksController.PutBook)
	router.DELETE("/api/v1/books/:bookid", booksController.DeleteBook)
	router.GET("/api/v1/books/:bookid/chapters", booksController.GetChapters)
	router.POST("/api/v1/books/:bookid/chapters", booksController.PostChapters)
	router.GET("/api/v1/books/:bookid/chapters/:chapterid", booksController.GetChapter)
	router.PUT("/api/v1/books/:bookid/chapters/:chapterid", booksController.PutChapter)
	router.DELETE("/api/v1/books/:bookid/chapters/:chapterid", booksController.DeleteChapter)
	return &Router{router}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
