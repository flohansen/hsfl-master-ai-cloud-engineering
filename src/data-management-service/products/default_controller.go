package products

import (
	"net/http"
)

type DefaultController struct {
	productRepository Repository
}

func NewDefaultController(productRepository Repository) *DefaultController {
	return &DefaultController{productRepository}
}

func (controller DefaultController) GetProducts(writer http.ResponseWriter, request *http.Request) {
	/*writer.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(writer).Encode(controller.productRepository.FindAll())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}*/
}

func (DefaultController) PostProducts(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (DefaultController) GetProduct(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (DefaultController) PutProduct(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (DefaultController) DeleteProduct(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}
