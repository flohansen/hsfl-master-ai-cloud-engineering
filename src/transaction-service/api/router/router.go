package router

import (
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/transactions"
	"net/http"
)

type Router struct {
	router http.Handler
}

func New(
	transactionController transactions.Controller,
) *Router {
	router := router.New()

	/*router.GET("/api/v1/books", booksController.GetBooks)
	router.POST("/api/v1/books", booksController.PostBooks)
	router.GET("/api/v1/books/:bookid", booksController.GetBook)
	router.PUT("/api/v1/books/:bookid", booksController.PutBook)
	router.DELETE("/api/v1/books/:bookid", booksController.DeleteBook)

	router.GET("/api/v1/books/:bookid/chapters", booksController.GetChapters)
	router.POST("/api/v1/books/:bookid/chapters", booksController.PostChapters)
	router.GET("/api/v1/books/:bookid/chapters/:chapterid", booksController.GetChapter)
	router.PUT("/api/v1/books/:bookid/chapters/:chapterid", booksController.PutChapter)
	router.DELETE("/api/v1/books/:bookid/chapters/:chapterid", booksController.DeleteChapter)
	*/
	return &Router{router}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
