package prices

import (
	"encoding/json"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/data-management-service/prices/model"
	"net/http"
	"strconv"
)

type DefaultController struct {
	priceRepository Repository
}

func NewDefaultController(priceRepository Repository) *DefaultController {
	return &DefaultController{priceRepository}
}

func (controller DefaultController) PostPrice(writer http.ResponseWriter, request *http.Request) {
	var requestData JsonFormatCreatePriceRequest
	if err := json.NewDecoder(request.Body).Decode(&requestData); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := controller.priceRepository.Create(&model.Price{
		UserId:    requestData.UserId,
		ProductId: requestData.ProductId,
		Price:     requestData.Price,
	}); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func (controller DefaultController) GetPrice(writer http.ResponseWriter, request *http.Request) {
	userId, err := strconv.ParseUint(request.Context().Value("userId").(string), 10, 64)
	productId, err := strconv.ParseUint(request.Context().Value("productId").(string), 10, 64)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	value, err := controller.priceRepository.FindByIds(productId, userId)
	if err != nil {
		if err.Error() == ErrorPriceNotFound {
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
	userId, err := strconv.ParseUint(request.Context().Value("userId").(string), 10, 64)
	productId, err := strconv.ParseUint(request.Context().Value("productId").(string), 10, 64)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	var requestData JsonFormatUpdatePriceRequest
	if err := json.NewDecoder(request.Body).Decode(&requestData); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := controller.priceRepository.Update(&model.Price{
		UserId:    userId,
		ProductId: productId,
		Price:     requestData.Price,
	}); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (controller DefaultController) DeleteProduct(writer http.ResponseWriter, request *http.Request) {
	userId, err := strconv.ParseUint(request.Context().Value("userId").(string), 10, 64)
	productId, err := strconv.ParseUint(request.Context().Value("productId").(string), 10, 64)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := controller.priceRepository.Delete(&model.Price{UserId: userId, ProductId: productId}); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
