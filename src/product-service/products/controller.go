package products

import "net/http"

type JsonFormatUpdateProductRequest struct {
	Description string `json:"description,omitempty"`
	Ean         uint64 `json:"ean,omitempty"`
}

type JsonFormatCreateProductRequest struct {
	Description string `json:"description,omitempty"`
	Ean         uint64 `json:"ean,omitempty"`
}

type Controller interface {
	GetProducts(http.ResponseWriter, *http.Request)
	GetProductByEan(http.ResponseWriter, *http.Request)
	GetProductById(http.ResponseWriter, *http.Request)
	PostProduct(http.ResponseWriter, *http.Request)
	PutProduct(http.ResponseWriter, *http.Request)
	DeleteProduct(http.ResponseWriter, *http.Request)
}
