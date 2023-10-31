package userShoppingListEntry

import "net/http"

type JsonFormatUpdateListRequest struct {
	ShoppingListId uint64 `json:"shoppingListId,omitempty"`
	ProductId      uint64 `json:"productId,omitempty"`
	Count          uint16 `json:"count,omitempty"`
	Note           string `json:"note,omitempty"`
	Checked        bool   `json:"checked,omitempty"`
}

type Controller interface {
	GetList(http.ResponseWriter, *http.Request)
	GetLists(http.ResponseWriter, *http.Request)
	PostList(http.ResponseWriter, *http.Request)
	PutList(http.ResponseWriter, *http.Request)
	DeleteList(http.ResponseWriter, *http.Request)
}
