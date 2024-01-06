package router

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router/middleware/test"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices"
	pricesMock "hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices/_mock"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products"
	productsMock "hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products/_mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouter(t *testing.T) {
	var mockProductsController = productsMock.NewMockController(t)
	var mockPricesController = pricesMock.NewMockController(t)

	authMiddleware := test.CreateEmptyMiddleware()
	var productsController products.Controller = mockProductsController
	var pricesController prices.Controller = mockPricesController

	router := New(&productsController, &pricesController, authMiddleware)

	t.Run("/api/v1/products", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/product/", nil)

			mockProductsController.EXPECT().GetProducts(w, r).Run(
				func(_a0 http.ResponseWriter, _a1 *http.Request) {
					_a0.WriteHeader(http.StatusOK)
				})

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call POST handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			jsonRequest := `{"id": 3, "description": "Test Product", "ean": 12345}`
			r := httptest.NewRequest("POST", "/api/v1/product/", strings.NewReader(jsonRequest))

			mockProductsController.EXPECT().PostProduct(w, r).Run(
				func(_a0 http.ResponseWriter, _a1 *http.Request) {
					_a0.WriteHeader(http.StatusOK)
				})

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/api/v1/product/:productid", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/product/1", nil)

			mockProductsController.EXPECT().GetProductById(w, mock.Anything).Run(
				func(_a0 http.ResponseWriter, _a1 *http.Request) {
					_a0.WriteHeader(http.StatusOK)
				})

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call PUT handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			jsonRequest := `{"description":"New Description","ean":4014819040771}`
			r := httptest.NewRequest("PUT", "/api/v1/product/1", strings.NewReader(jsonRequest))

			mockProductsController.EXPECT().PutProduct(w, mock.Anything).Run(
				func(_a0 http.ResponseWriter, _a1 *http.Request) {
					_a0.WriteHeader(http.StatusOK)
				})

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call DELETE handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/product/1", nil)

			mockProductsController.EXPECT().DeleteProduct(w, mock.Anything).Run(
				func(_a0 http.ResponseWriter, _a1 *http.Request) {
					_a0.WriteHeader(http.StatusOK)
				})

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/api/v1/price", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/price/", nil)

			mockPricesController.EXPECT().GetPrices(w, r).Run(
				func(_a0 http.ResponseWriter, _a1 *http.Request) {
					_a0.WriteHeader(http.StatusOK)
				})

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call POST handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			jsonRequest := `{"price": 0.99}`
			r := httptest.NewRequest("POST", "/api/v1/price/3/3", strings.NewReader(jsonRequest))

			mockPricesController.EXPECT().PostPrice(w, mock.Anything).Run(
				func(_a0 http.ResponseWriter, _a1 *http.Request) {
					_a0.WriteHeader(http.StatusOK)
				})

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/api/v1/price/user/:userId", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/price/user/1", nil)

			mockPricesController.EXPECT().GetPricesByUser(w, mock.Anything).Run(
				func(_a0 http.ResponseWriter, _a1 *http.Request) {
					_a0.WriteHeader(http.StatusOK)
				})

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/api/v1/price/:productId/:userId", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/price/1/1", nil)

			mockPricesController.EXPECT().GetPrice(w, mock.Anything).Run(
				func(_a0 http.ResponseWriter, _a1 *http.Request) {
					_a0.WriteHeader(http.StatusOK)
				})

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call PUT handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			jsonRequest := `{"userId": 1, "productId": 1, "price": 10.99}`
			r := httptest.NewRequest("PUT", "/api/v1/price/1/1", strings.NewReader(jsonRequest))

			mockPricesController.EXPECT().PutPrice(w, mock.Anything).Run(
				func(_a0 http.ResponseWriter, _a1 *http.Request) {
					_a0.WriteHeader(http.StatusOK)
				})

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call DELETE handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/price/1/1", nil)

			mockPricesController.EXPECT().DeletePrice(w, mock.Anything).Run(
				func(_a0 http.ResponseWriter, _a1 *http.Request) {
					_a0.WriteHeader(http.StatusOK)
				})

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
}
