package service

import "net/http"

type Controller interface {
	GetShoppingList(http.ResponseWriter, *http.Request)
	GetAdminProducts(http.ResponseWriter, *http.Request)
	GetMerchantProducts(http.ResponseWriter, *http.Request)
	GetProductCatalogue(http.ResponseWriter, *http.Request)
}
