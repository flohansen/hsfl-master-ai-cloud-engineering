package products

import (
	"encoding/json"
	"golang.org/x/sync/singleflight"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router/middleware/auth"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products/model"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products/utils"
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

func (controller *coalescingController) GetProducts(writer http.ResponseWriter, request *http.Request) {
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

func (controller *coalescingController) GetProductById(writer http.ResponseWriter, request *http.Request) {
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

func (controller *coalescingController) GetProductByEan(writer http.ResponseWriter, request *http.Request) {
	productEan, exists := request.Context().Value("productEan").(string)
	if !exists || utils.ValidateEAN(productEan) == false {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	msg, err, _ := controller.group.Do("get_ean_"+productEan, func() (interface{}, error) {
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

func (controller *coalescingController) PostProduct(writer http.ResponseWriter, request *http.Request) {
	authUserRole, _ := request.Context().Value("auth_userRole").(int64)

	if authUserRole == auth.Administrator || authUserRole == auth.Merchant {
		var requestData JsonFormatCreateProductRequest
		if err := json.NewDecoder(request.Body).Decode(&requestData); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		if utils.ValidateEAN(requestData.Ean) == false {
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
	} else {
		writer.WriteHeader(http.StatusUnauthorized)
	}
}

func (controller *coalescingController) PutProduct(writer http.ResponseWriter, request *http.Request) {
	authUserRole, _ := request.Context().Value("auth_userRole").(int64)

	if authUserRole == auth.Administrator || authUserRole == auth.Merchant {
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

		if utils.ValidateEAN(requestData.Ean) == false {
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
	} else {
		writer.WriteHeader(http.StatusUnauthorized)
	}
}

func (controller *coalescingController) DeleteProduct(writer http.ResponseWriter, request *http.Request) {
	productId, err := strconv.ParseUint(request.Context().Value("productId").(string), 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	authUserRole, _ := request.Context().Value("auth_userRole").(int64)

	if authUserRole == auth.Administrator {
		if err := controller.productRepository.Delete(&model.Product{Id: productId}); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		writer.WriteHeader(http.StatusUnauthorized)
	}
}
