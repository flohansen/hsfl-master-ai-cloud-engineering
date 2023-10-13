package service

import "net/http"

type DefaultController struct{}

func NewDefaultController() *DefaultController {
	return &DefaultController{}
}

func (d DefaultController) GetShoppingList(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Shopping List"))
	if err != nil {
		return
	}
}

func (d DefaultController) GetAdminProducts(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Admin Products"))
	if err != nil {
		return
	}
}

func (d DefaultController) GetMerchantProducts(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Merchant Products"))
	if err != nil {
		return
	}
}

func (d DefaultController) GetProductCatalogue(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Product Catalogue"))
	if err != nil {
		return
	}
}
