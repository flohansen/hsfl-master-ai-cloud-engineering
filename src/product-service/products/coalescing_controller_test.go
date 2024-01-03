package products

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/singleflight"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products/model"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
)

func TestNewCoalescingController(t *testing.T) {
	demoRepo := NewDemoRepository()

	controller := NewCoalescingController(demoRepo)

	assert.NotNil(t, controller)
	assert.Equal(t, demoRepo, controller.productRepository)
	assert.IsType(t, &singleflight.Group{}, controller.group)
}

func TestCoalescingController_DeleteProduct(t *testing.T) {
	type fields struct {
		productRepository Repository
	}
	type args struct {
		writer  *httptest.ResponseRecorder
		request *http.Request
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
	}{
		{
			name: "Successfully delete existing product (expect 200)",
			fields: fields{
				productRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest("DELETE", "/api/v1/product/1", nil)
					request = request.WithContext(context.WithValue(request.Context(), "productId", "1"))
					return request
				}(),
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Bad non-numeric request (expect 400)",
			fields: fields{
				productRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest("DELETE", "/api/v1/product/abc", nil)
					request = request.WithContext(context.WithValue(request.Context(), "productId", "abc"))
					return request
				}(),
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Unknown product to delete (expect 500)",
			fields: fields{
				productRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest("DELETE", "/api/v1/product/5", nil)
					request = request.WithContext(context.WithValue(request.Context(), "productId", "5"))
					return request
				}(),
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := coalescingController{
				productRepository: tt.fields.productRepository,
			}
			controller.DeleteProduct(tt.args.writer, tt.args.request)
			if tt.args.writer.Code != tt.wantStatus {
				t.Errorf("Expected status code %d, got %d", tt.wantStatus, tt.args.writer.Code)
			}
		})
	}
}

func TestCoalescingController_GetProductById(t *testing.T) {
	type fields struct {
		productRepository Repository
	}
	type args struct {
		writer  *httptest.ResponseRecorder
		request *http.Request
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
	}{
		{
			name: "Bad non-numeric request (expect 400)",
			fields: fields{
				productRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest("GET", "/api/v1/product/abc", nil)
					request = request.WithContext(context.WithValue(request.Context(), "productId", "abc"))
					return request
				}(),
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Unknown product (expect 404)",
			fields: fields{
				productRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest("GET", "/api/v1/product/4", nil)
					request = request.WithContext(context.WithValue(request.Context(), "productId", "4"))
					return request
				}(),
			},
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := NewCoalescingController(tt.fields.productRepository)
			controller.GetProductById(tt.args.writer, tt.args.request)
			if tt.args.writer.Code != tt.wantStatus {
				t.Errorf("Expected status code %d, got %d", tt.wantStatus, tt.args.writer.Code)
			}
		})
	}

	t.Run("Successfully get existing product (expect 200 and product)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/product/1", nil)
		request = request.WithContext(context.WithValue(request.Context(), "productId", "1"))

		controller := NewCoalescingController(GenerateExampleDemoRepository())

		// when
		controller.GetProductById(writer, request)

		// then
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

		if response.Ean != 4014819040771 {
			t.Errorf("Expected ean of product %d, got %d", 4014819040771, response.Ean)
		}

	})
}

func TestCoalescingController_GetProductByEan(t *testing.T) {
	t.Run("Bad non-numeric request (expect 400)", func(t *testing.T) {
		controller := NewCoalescingController(GenerateExampleDemoRepository())

		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/products/ean?ean=abc", nil)
		request = request.WithContext(context.WithValue(request.Context(), "productEan", "abc"))

		// Test request
		controller.GetProductByEan(writer, request)

		if writer.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, writer.Code)
		}
	})

	t.Run("Unknown product (expect 404)", func(t *testing.T) {
		controller := NewCoalescingController(GenerateExampleDemoRepository())

		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/products/ean?ean=123", nil)
		request = request.WithContext(context.WithValue(request.Context(), "productEan", "123"))

		// Test request
		controller.GetProductByEan(writer, request)

		if writer.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, writer.Code)
		}
	})

	t.Run("Should return products by EAN", func(t *testing.T) {
		controller := NewCoalescingController(GenerateExampleDemoRepository())

		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/products/ean?ean=4014819040771", nil)
		request = request.WithContext(context.WithValue(request.Context(), "productEan", "4014819040771"))

		// Test request
		controller.GetProductByEan(writer, request)

		if writer.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, writer.Code)
		}

		res := writer.Result()
		var response model.Product
		err := json.NewDecoder(res.Body).Decode(&response)

		if err != nil {
			t.Error(err)
		}

		// Add assertions based on your expected response.
		// For example, check the length, content, etc.
		// ...

	})
}

func TestCoalescingController_GetProducts(t *testing.T) {
	t.Run("should return all products", func(t *testing.T) {
		controller := NewCoalescingController(GenerateExampleDemoRepository())

		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/product", nil)

		// Test request
		controller.GetProducts(writer, request)

		res := writer.Result()
		var response []model.Product
		err := json.NewDecoder(res.Body).Decode(&response)

		if err != nil {
			t.Error(err)
		}

		if writer.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, writer.Code)
		}

		if writer.Header().Get("Content-Type") != "application/json" {
			t.Errorf("Expected content type %s, got %s",
				"application/json", writer.Header().Get("Content-Type"))
		}

		products := GenerateExampleProductSlice()

		sort.Slice(response, func(i, j int) bool {
			return response[i].Id < response[j].Id
		})

		if len(response) != len(products) {
			t.Errorf("Expected count of product is %d, got %d",
				2, len(response))
		}

		for i, product := range products {
			if product.Id != response[i].Id {
				t.Errorf("Expected id of product %d, got %d", product.Id, response[i].Id)
			}

			if product.Description != response[i].Description {
				t.Errorf("Expected description of product %s, got %s", product.Description, response[i].Description)
			}

			if product.Ean != response[i].Ean {
				t.Errorf("Expected ean of product %d, got %d", product.Ean, response[i].Ean)
			}
		}

	})
}

func TestCoalescingController_PostProduct(t *testing.T) {
	type fields struct {
		productRepository Repository
	}
	type args struct {
		writer  *httptest.ResponseRecorder
		request *http.Request
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "Valid Product",
			fields: fields{
				productRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: httptest.NewRequest(
					"POST",
					"/api/v1/product",
					strings.NewReader(`{"id": 3, "description": "Test Product", "ean": 12345}`),
				),
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "",
		},
		{
			name: "Valid Product (Partly Fields)",
			fields: fields{
				productRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: httptest.NewRequest(
					"POST",
					"/api/v1/product",
					strings.NewReader(`{"description": "Incomplete Product"}`),
				),
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "",
		},
		{
			name: "Malformed JSON",
			fields: fields{
				productRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: httptest.NewRequest(
					"POST",
					"/api/v1/product",
					strings.NewReader(`{"description": "Incomplete Product"`),
				),
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "",
		},
		{
			name: "Invalid product, incorrect Type for EAN (Non-numeric)",
			fields: fields{
				productRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: httptest.NewRequest(
					"POST",
					"/api/v1/product",
					strings.NewReader(`{"id": 3, "description": "Invalid EAN", "ean": "abc"}`),
				),
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := coalescingController{
				productRepository: tt.fields.productRepository,
			}
			controller.PostProduct(tt.args.writer, tt.args.request)

			// You can then assert the response status and content, and check against your expectations.
			if tt.args.writer.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, but got %d", tt.expectedStatus, tt.args.writer.Code)
			}

			if tt.expectedResponse != "" {
				actualResponse := tt.args.writer.Body.String()
				if actualResponse != tt.expectedResponse {
					t.Errorf("Expected response: %s, but got: %s", tt.expectedResponse, actualResponse)
				}
			}
		})
	}
}

func TestCoalescingController_PutProduct(t *testing.T) {
	type fields struct {
		productRepository Repository
	}
	type args struct {
		writer  *httptest.ResponseRecorder
		request *http.Request
	}

	tests := []struct {
		name             string
		fields           fields
		args             args
		expectedStatus   int
		expectedResponse string // If you want to check the response content
	}{
		{
			name: "Valid Update",
			fields: fields{
				productRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/product/1",
						strings.NewReader(`{"id": 1, "description": "Updated Product", "ean": 54321}`))
					request = request.WithContext(context.WithValue(request.Context(), "productId", "1"))
					return request
				}(),
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "",
		},
		{
			name: "Valid Update (Partly Fields)",
			fields: fields{
				productRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/product/2",
						strings.NewReader(`{"description": "Incomplete Update"}`))
					request = request.WithContext(context.WithValue(request.Context(), "productId", "2"))
					return request
				}(),
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "",
		},
		{
			name: "Malformed JSON",
			fields: fields{
				productRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/product/2",
						strings.NewReader(`{"description": "Incomplete Update"`))
					request = request.WithContext(context.WithValue(request.Context(), "productId", "2"))
					return request
				}(),
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "",
		},
		{
			name:   "Incorrect Type for EAN (Non-numeric)",
			fields: fields{
				// Set up your repository mock or test double here if needed
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/product/2",
						strings.NewReader(`{"ean": "Wrong Type"`))
					request = request.WithContext(context.WithValue(request.Context(), "productId", "2"))
					return request
				}(),
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := coalescingController{
				productRepository: tt.fields.productRepository,
			}
			controller.PutProduct(tt.args.writer, tt.args.request)

			// You can then assert the response status and content, and check against your expectations.
			if tt.args.writer.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, but got %d", tt.expectedStatus, tt.args.writer.Code)
			}

			if tt.expectedResponse != "" {
				actualResponse := tt.args.writer.Body.String()
				if actualResponse != tt.expectedResponse {
					t.Errorf("Expected response: %s, but got: %s", tt.expectedResponse, actualResponse)
				}
			}
		})
	}
}
