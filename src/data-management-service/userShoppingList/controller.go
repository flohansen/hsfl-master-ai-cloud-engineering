package userShoppingList

import "net/http"

type Controller interface {
	GetList(http.ResponseWriter, *http.Request)
	PutList(http.ResponseWriter, *http.Request)
	DeleteList(http.ResponseWriter, *http.Request)
}
