package prices

import (
	"encoding/json"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices/model"
	"net/http"
	"strconv"
)

type defaultController struct {
	priceRepository Repository
}

func NewDefaultController(priceRepository Repository) *defaultController {
	return &defaultController{priceRepository}
}

func (controller defaultController) GetPrices(writer http.ResponseWriter, request *http.Request) {
	values, err := controller.priceRepository.FindAll()
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

func (controller defaultController) GetPricesByUser(writer http.ResponseWriter, request *http.Request) {
	userId, err := strconv.ParseUint(request.Context().Value("userId").(string), 10, 64)

	if err != nil {
		http.Error(writer, "Invalid userId", http.StatusBadRequest)
		return
	}
	values, err := controller.priceRepository.FindAllByUser(userId)
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

func (controller defaultController) PostPrice(writer http.ResponseWriter, request *http.Request) {
	productId, productIdErr := strconv.ParseUint(request.Context().Value("productId").(string), 10, 64)
	userId, userIdErr := strconv.ParseUint(request.Context().Value("userId").(string), 10, 64)

	if productIdErr != nil || userIdErr != nil {
		http.Error(writer, "Invalid listId or productId", http.StatusBadRequest)
		return
	}

	var requestData JsonFormatCreatePriceRequest
	if err := json.NewDecoder(request.Body).Decode(&requestData); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := controller.priceRepository.Create(&model.Price{
		ProductId: productId,
		UserId:    userId,
		Price:     requestData.Price,
	}); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func (controller defaultController) GetPrice(writer http.ResponseWriter, request *http.Request) {
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

func (controller defaultController) PutPrice(writer http.ResponseWriter, request *http.Request) {
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

func (controller defaultController) DeletePrice(writer http.ResponseWriter, request *http.Request) {
	userId, err := strconv.ParseUint(request.Context().Value("userId").(string), 10, 64)
	productId, err := strconv.ParseUint(request.Context().Value("productId").(string), 10, 64)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := controller.priceRepository.Delete(&model.Price{ProductId: productId, UserId: userId}); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
