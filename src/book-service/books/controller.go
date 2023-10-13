package books

import "net/http"

type Controller interface {
	GetBooks(http.ResponseWriter, *http.Request)
	GetBook(http.ResponseWriter, *http.Request)
	GetChapters(http.ResponseWriter, *http.Request)
	GetChapter(http.ResponseWriter, *http.Request)
	PostBooks(http.ResponseWriter, *http.Request)
	PutBook(http.ResponseWriter, *http.Request)
	DeleteBook(http.ResponseWriter, *http.Request)
	PostChapters(http.ResponseWriter, *http.Request)
	PutChapter(http.ResponseWriter, *http.Request)
	DeleteChapter(http.ResponseWriter, *http.Request)
}
