package products

import (
	"encoding/json"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/data-management-service/products/model"
	"net/http"
	"strconv"
)

type DefaultController struct {
	productRepository Repository
}

func NewDefaultController(productRepository Repository) *DefaultController {
	return &DefaultController{productRepository}
}

func (controller DefaultController) GetProducts(writer http.ResponseWriter, request *http.Request) {
	values, err := controller.productRepository.FindAll()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(values)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (controller DefaultController) PostProduct(writer http.ResponseWriter, request *http.Request) {
	var requestData JsonFormatCreateProductRequest
	if err := json.NewDecoder(request.Body).Decode(&requestData); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := controller.productRepository.Create(&model.Product{
		Description: requestData.Description,
		Ean:         requestData.Ean,
	}); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (controller DefaultController) GetProduct(writer http.ResponseWriter, request *http.Request) {
	productId, err := strconv.ParseUint(request.Context().Value("productId").(string), 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	value, err := controller.productRepository.FindById(productId)
	if err != nil {
		if err.Error() == ErrorProductNotFound {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(value)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (controller DefaultController) PutProduct(writer http.ResponseWriter, request *http.Request) {
	productId, err := strconv.ParseUint(request.Context().Value("productId").(string), 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	var requestData JsonFormatUpdateProductRequest
	if err := json.NewDecoder(request.Body).Decode(&requestData); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := controller.productRepository.Update(&model.Product{
		Id:          productId,
		Description: requestData.Description,
		Ean:         requestData.Ean,
	}); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (controller DefaultController) DeleteProduct(writer http.ResponseWriter, request *http.Request) {
	productId, err := strconv.ParseUint(request.Context().Value("productId").(string), 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := controller.productRepository.Delete(&model.Product{Id: productId}); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
