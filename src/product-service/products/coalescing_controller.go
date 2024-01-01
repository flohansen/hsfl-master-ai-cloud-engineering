package products

import (
	"encoding/json"
	"golang.org/x/sync/singleflight"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products/model"
	"net/http"
	"strconv"
)

type coalescingController struct {
	productRepository Repository
	group             *singleflight.Group
}

func NewCoalescingController(productRepository Repository) *coalescingController {
	return &coalescingController{
		productRepository,
		&singleflight.Group{},
	}
}

func (controller coalescingController) GetProducts(writer http.ResponseWriter, request *http.Request) {
	msg, err, _ := controller.group.Do("get-all", func() (interface{}, error) {
		return controller.productRepository.FindAll()
	})

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(msg)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (controller coalescingController) PostProduct(writer http.ResponseWriter, request *http.Request) {
	var requestData JsonFormatCreateProductRequest
	if err := json.NewDecoder(request.Body).Decode(&requestData); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := controller.productRepository.Create(&model.Product{
		Description: requestData.Description,
		Ean:         requestData.Ean,
	})

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(product)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (controller coalescingController) GetProductById(writer http.ResponseWriter, request *http.Request) {
	productIdAttribute := request.Context().Value("productId").(string)

	msg, err, _ := controller.group.Do("get_id_"+productIdAttribute, func() (interface{}, error) {
		productId, err := strconv.ParseUint(productIdAttribute, 10, 64)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return nil, err
		}

		value, err := controller.productRepository.FindById(productId)
		if err != nil {
			if err.Error() == ErrorProductNotFound {
				http.Error(writer, err.Error(), http.StatusNotFound)
				return nil, err
			}
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}

		return value, nil
	})

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(msg)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (controller coalescingController) GetProductByEan(writer http.ResponseWriter, request *http.Request) {
	productEanAttribute := request.Context().Value("productEan").(string)

	msg, err, _ := controller.group.Do("get_ean_"+productEanAttribute, func() (interface{}, error) {
		productEan, err := strconv.ParseUint(productEanAttribute, 10, 64)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return nil, err
		}

		value, err := controller.productRepository.FindByEan(productEan)
		if err != nil {
			if err.Error() == ErrorProductNotFound {
				http.Error(writer, err.Error(), http.StatusNotFound)
				return nil, err
			}
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}

		return value, nil
	})

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(msg)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (controller coalescingController) PutProduct(writer http.ResponseWriter, request *http.Request) {
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

func (controller coalescingController) DeleteProduct(writer http.ResponseWriter, request *http.Request) {
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
