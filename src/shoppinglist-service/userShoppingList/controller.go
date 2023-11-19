package userShoppingList

import "net/http"

type JsonFormatUpdateListRequest struct {
	UserId      uint64 `json:"userId,omitempty"`
	Description string `json:"description,omitempty"`
	Checked     bool   `json:"checked,omitempty"`
}

type JsonFormatCreateListRequest struct {
	UserId      uint64 `json:"userId,omitempty"`
	Description string `json:"description,omitempty"`
}

type Controller interface {
	GetList(http.ResponseWriter, *http.Request)
	GetLists(http.ResponseWriter, *http.Request)
	PostList(http.ResponseWriter, *http.Request)
	PutList(http.ResponseWriter, *http.Request)
	DeleteList(http.ResponseWriter, *http.Request)
}
