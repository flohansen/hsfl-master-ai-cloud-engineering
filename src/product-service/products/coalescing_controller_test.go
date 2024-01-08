package products

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/singleflight"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router/middleware/auth"
	products "hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products/_mock"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewCoalescingController(t *testing.T) {
	var mockProductRepository = products.NewMockRepository(t)
	var productRepository Repository = mockProductRepository

	productController := NewCoalescingController(productRepository)

	assert.NotNil(t, productController)
	assert.Equal(t, productRepository, productController.productRepository)
	assert.IsType(t, &singleflight.Group{}, productController.group)
}

func TestCoalescingController_GetProducts(t *testing.T) {
	var mockProductRepository = products.NewMockRepository(t)
	var productRepository Repository = mockProductRepository

	productController := NewCoalescingController(productRepository)

	t.Run("Successfully return all products (expect 200 and products)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/product", nil)

		mockProductRepository.EXPECT().FindAll().Return([]*model.Product{
			{
				Id:          1,
				Description: "Strauchtomaten",
				Ean:         "4014819040771",
			},
			{
				Id:          2,
				Description: "Lauchzwiebeln",
				Ean:         "5001819040871",
			},
		}, nil)

		productController.GetProducts(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)

		if writer.Header().Get("Content-Type") != "application/json" {
			t.Errorf("Expected content type %s, got %s",
				"application/json", writer.Header().Get("Content-Type"))
		}

		res := writer.Result()
		var response []model.Product
		err := json.NewDecoder(res.Body).Decode(&response)

		if err != nil {
			t.Error(err)
		}

		if len(response) != 2 {
			t.Errorf("Expected count of product is %d, got %d",
				2, len(response))
		}
	})
}

func TestCoalescingController_GetProductById(t *testing.T) {
	var mockProductRepository = products.NewMockRepository(t)
	var productRepository Repository = mockProductRepository

	productController := NewCoalescingController(productRepository)

	t.Run("Bad non-numeric request (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/product/abc", nil)
		ctx := context.WithValue(request.Context(), "productId", "abc")
		request = request.WithContext(ctx)

		productController.GetProductById(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Unknown product (expect 404)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/product/4", nil)
		ctx := context.WithValue(request.Context(), "auth_userId", uint64(1))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Customer)
		ctx = context.WithValue(ctx, "productId", "4")
		request = request.WithContext(ctx)

		mockProductRepository.EXPECT().FindById(uint64(4)).Return(nil, errors.New(ErrorProductNotFound))

		productController.GetProductById(writer, request)

		assert.Equal(t, http.StatusNotFound, writer.Code)
	})

	t.Run("Successfully get existing product (expect 200 and product)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/product/1", nil)
		ctx := context.WithValue(request.Context(), "productId", "1")
		request = request.WithContext(ctx)

		mockProductRepository.EXPECT().FindById(uint64(1)).Return(&model.Product{
			Id:          1,
			Description: "Strauchtomaten",
			Ean:         "4014819040771",
		}, nil)

		productController.GetProductById(writer, request)

		if writer.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, writer.Code)
		}

		if writer.Header().Get("Content-Type") != "application/json" {
			t.Errorf("Expected content type %s, got %s",
				"application/json", writer.Header().Get("Content-Type"))
		}

		result := writer.Result()
		var response model.Product
		err := json.NewDecoder(result.Body).Decode(&response)
		if err != nil {
			t.Fatal(err.Error())
		}

		if response.Id != 1 {
			t.Errorf("Expected id of product %d, got %d", 1, response.Id)
		}

		if response.Description != "Strauchtomaten" {
			t.Errorf("Expected description of product %s, got %s", "Strauchtomaten", response.Description)
		}

		if response.Ean != "4014819040771" {
			t.Errorf("Expected ean of product %s, got %s", "4014819040771", response.Ean)
		}

	})
}

func TestCoalescingController_GetProductByEan(t *testing.T) {
	var mockProductRepository = products.NewMockRepository(t)
	var productRepository Repository = mockProductRepository

	productController := NewCoalescingController(productRepository)

	t.Run("Invalid product ean (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/products/ean?ean=123", nil)
		request = request.WithContext(context.WithValue(request.Context(), "productEan", "123"))

		productController.GetProductByEan(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Unknown product (expect 404)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/products/ean?ean=12345670", nil)
		request = request.WithContext(context.WithValue(request.Context(), "productEan", "12345670"))

		mockProductRepository.EXPECT().FindByEan("12345670").Return(nil, errors.New(ErrorProductNotFound))

		productController.GetProductByEan(writer, request)

		assert.Equal(t, http.StatusNotFound, writer.Code)
	})

	t.Run("Should return product by EAN (expect 200 and product)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/products/ean?ean=4014819040771", nil)
		request = request.WithContext(context.WithValue(request.Context(), "productEan", "4014819040771"))

		mockProductRepository.EXPECT().FindByEan("4014819040771").Return(&model.Product{
			Id:          1,
			Description: "Strauchtomaten",
			Ean:         "4014819040771",
		}, nil)

		productController.GetProductByEan(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)

		res := writer.Result()
		var response model.Product
		err := json.NewDecoder(res.Body).Decode(&response)
		if err != nil {
			t.Error(err)
		}

		if response.Id != 1 {
			t.Errorf("Expected id of product %d, got %d", 1, response.Id)
		}

		if response.Description != "Strauchtomaten" {
			t.Errorf("Expected description of product %s, got %s", "Strauchtomaten", response.Description)
		}

		if response.Ean != "4014819040771" {
			t.Errorf("Expected ean of product %s, got %s", "4014819040771", response.Ean)
		}
	})
}

func TestCoalescingController_PostProduct(t *testing.T) {
	var mockProductRepository = products.NewMockRepository(t)
	var productRepository Repository = mockProductRepository

	productController := NewCoalescingController(productRepository)

	t.Run("Unauthorized (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/v1/product",
			strings.NewReader(`{"id": 3, "description": "Test Product", "ean": "12345"}`))

		productController.PostProduct(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Invalid create (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/v1/product",
			strings.NewReader(`{"description": "Test Product", "ean": "12345"}`))
		ctx := context.WithValue(request.Context(), "auth_userRole", int64(1))
		ctx = context.WithValue(ctx, "auth_userId", uint(1))
		request = request.WithContext(ctx)

		productController.PostProduct(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Valid create (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/v1/product",
			strings.NewReader(`{"description": "Test Product", "ean": "12345670"}`))
		ctx := context.WithValue(request.Context(), "auth_userRole", auth.Merchant)
		ctx = context.WithValue(ctx, "auth_userId", uint(1))
		request = request.WithContext(ctx)

		mockProductRepository.EXPECT().Create(&model.Product{
			Description: "Test Product",
			Ean:         "12345670",
		}).Return(&model.Product{
			Id:          3,
			Description: "Test Product",
			Ean:         "12345670",
		}, nil)

		productController.PostProduct(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Valid create ad admin(expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/v1/product",
			strings.NewReader(`{"description": "Test Product", "ean": "12345670"}`))
		ctx := context.WithValue(request.Context(), "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		mockProductRepository.EXPECT().Create(&model.Product{
			Description: "Test Product",
			Ean:         "12345670",
		}).Return(&model.Product{
			Id:          3,
			Description: "Test Product",
			Ean:         "12345670",
		}, nil)

		productController.PostProduct(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Malformed JSON (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/v1/product",
			strings.NewReader(`{"description": "Test Product", "ean": "12345670"`))
		ctx := context.WithValue(request.Context(), "auth_userRole", auth.Merchant)
		ctx = context.WithValue(ctx, "auth_userId", uint(1))
		request = request.WithContext(ctx)

		productController.PostProduct(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})
}

func TestCoalescingController_PutProduct(t *testing.T) {
	var mockProductRepository = products.NewMockRepository(t)
	var productRepository Repository = mockProductRepository

	productController := NewCoalescingController(productRepository)

	t.Run("Unauthorized (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/product/1",
			strings.NewReader(`{"description": "Updated Product", "ean": "54321"}`))
		ctx := context.WithValue(request.Context(), "productId", "1")
		request = request.WithContext(ctx)

		productController.PutProduct(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Valid update (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/product/1",
			strings.NewReader(`{"description": "Updated Product", "ean": "4001686323397"}`))
		ctx := context.WithValue(request.Context(), "auth_userRole", auth.Merchant)
		ctx = context.WithValue(ctx, "auth_userId", uint(1))
		ctx = context.WithValue(ctx, "productId", "1")
		request = request.WithContext(ctx)

		mockProductRepository.EXPECT().Update(&model.Product{
			Id:          1,
			Description: "Updated Product",
			Ean:         "4001686323397",
		}).Return(&model.Product{
			Id:          1,
			Description: "Updated Product",
			Ean:         "4001686323397",
		}, nil)

		productController.PutProduct(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Invalid update (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/product/1",
			strings.NewReader(`{"description": "Suppe", "ean": "54321"}`))
		ctx := context.WithValue(request.Context(), "auth_userRole", auth.Merchant)
		ctx = context.WithValue(ctx, "auth_userId", uint(1))
		ctx = context.WithValue(ctx, "productId", "1")
		request = request.WithContext(ctx)

		productController.PutProduct(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Malformed JSON (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/product/1",
			strings.NewReader(`{"description": "Suppe", "ean": "12345679"`))
		ctx := context.WithValue(request.Context(), "auth_userRole", auth.Merchant)
		ctx = context.WithValue(ctx, "auth_userId", uint(1))
		ctx = context.WithValue(ctx, "productId", "1")
		request = request.WithContext(ctx)

		productController.PutProduct(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Bad type for EAN (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/product/1",
			strings.NewReader(`{"description": "Suppe", "ean": "test"}`))
		ctx := context.WithValue(request.Context(), "auth_userRole", auth.Merchant)
		ctx = context.WithValue(ctx, "auth_userId", uint(1))
		ctx = context.WithValue(ctx, "productId", "1")
		request = request.WithContext(ctx)

		productController.PutProduct(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Unknown product (expect 500)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/product/10",
			strings.NewReader(`{"description": "Updated Product", "ean": "4001686323397"}`))
		ctx := context.WithValue(request.Context(), "auth_userRole", auth.Merchant)
		ctx = context.WithValue(ctx, "auth_userId", uint(1))
		ctx = context.WithValue(ctx, "productId", "10")
		request = request.WithContext(ctx)

		mockProductRepository.EXPECT().Update(&model.Product{
			Id:          10,
			Description: "Updated Product",
			Ean:         "4001686323397",
		}).Return(nil, errors.New(ErrorProductNotFound))

		productController.PutProduct(writer, request)

		assert.Equal(t, http.StatusInternalServerError, writer.Code)
	})
}

func TestCoalescingController_DeleteProduct(t *testing.T) {
	var mockProductRepository = products.NewMockRepository(t)
	var productRepository Repository = mockProductRepository

	productController := NewCoalescingController(productRepository)

	t.Run("Unauthorized (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/product/1",
			strings.NewReader(`{"description": "Suppe", "ean": "54321"}`))
		ctx := context.WithValue(request.Context(), "productId", "1")
		request = request.WithContext(ctx)

		productController.DeleteProduct(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Valid delete as admin (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/product/1", nil)
		ctx := context.WithValue(request.Context(), "productId", "1")
		ctx = context.WithValue(ctx, "auth_userId", "2")
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		mockProductRepository.EXPECT().Delete(&model.Product{Id: 1}).Return(nil)

		productController.DeleteProduct(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Invalid delete, non-numeric request (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/product/abc", nil)
		ctx := context.WithValue(request.Context(), "productId", "abc")
		ctx = context.WithValue(ctx, "auth_userId", "2")
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		productController.DeleteProduct(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Unknown product (expect 500)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/product/10", nil)
		ctx := context.WithValue(request.Context(), "productId", "10")
		ctx = context.WithValue(ctx, "auth_userId", "2")
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		mockProductRepository.EXPECT().Delete(&model.Product{Id: 10}).Return(errors.New(ErrorProductNotFound))

		productController.DeleteProduct(writer, request)

		assert.Equal(t, http.StatusInternalServerError, writer.Code)
	})
}
