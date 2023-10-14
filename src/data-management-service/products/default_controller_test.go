package products

import (
	"context"
	"encoding/json"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/data-management-service/products/model"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewDefaultController(t *testing.T) {
	type args struct {
		productRepository Repository
	}
	tests := []struct {
		name string
		args args
		want *defaultController
	}{
		{
			name: "Test construction with DemoRepository",
			args: args{productRepository: NewDemoRepository()},
			want: &defaultController{productRepository: NewDemoRepository()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDefaultController(tt.args.productRepository); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDefaultController() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultController_DeleteProduct(t *testing.T) {
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
				productRepository: setupMockRepository(),
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
				productRepository: setupMockRepository(),
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
				productRepository: setupMockRepository(),
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
			controller := defaultController{
				productRepository: tt.fields.productRepository,
			}
			controller.DeleteProduct(tt.args.writer, tt.args.request)
			if tt.args.writer.Code != tt.wantStatus {
				t.Errorf("Expected status code %d, got %d", tt.wantStatus, tt.args.writer.Code)
			}
		})
	}
}

func TestDefaultController_GetProduct(t *testing.T) {
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
				productRepository: setupMockRepository(),
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
				productRepository: setupMockRepository(),
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
			controller := defaultController{
				productRepository: tt.fields.productRepository,
			}
			controller.GetProduct(tt.args.writer, tt.args.request)
			if tt.args.writer.Code != tt.wantStatus {
				t.Errorf("Expected status code %d, got %d", tt.wantStatus, tt.args.writer.Code)
			}
		})
	}

	t.Run("Successfully get existing product (expect 200 and product)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/product/1", nil)
		request = request.WithContext(context.WithValue(request.Context(), "productId", "1"))

		controller := defaultController{
			productRepository: setupMockRepository(),
		}

		// when
		controller.GetProduct(writer, request)

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

func TestDefaultController_GetProducts(t *testing.T) {
	type fields struct {
		productRepository Repository
	}
	type args struct {
		writer  http.ResponseWriter
		request *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := defaultController{
				productRepository: tt.fields.productRepository,
			}
			controller.GetProducts(tt.args.writer, tt.args.request)
		})
	}
}

func TestDefaultController_PostProduct(t *testing.T) {
	type fields struct {
		productRepository Repository
	}
	type args struct {
		writer  http.ResponseWriter
		request *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := defaultController{
				productRepository: tt.fields.productRepository,
			}
			controller.PostProduct(tt.args.writer, tt.args.request)
		})
	}
}

func TestDefaultController_PutProduct(t *testing.T) {
	type fields struct {
		productRepository Repository
	}
	type args struct {
		writer  http.ResponseWriter
		request *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := defaultController{
				productRepository: tt.fields.productRepository,
			}
			controller.PutProduct(tt.args.writer, tt.args.request)
		})
	}
}

func setupMockRepository() Repository {
	repository := NewDemoRepository()
	productSlice := []*model.Product{
		{
			Id:          1,
			Description: "Strauchtomaten",
			Ean:         4014819040771,
		},
		{
			Id:          2,
			Description: "Lauchzwiebeln",
			Ean:         5001819040871,
		},
	}
	for _, product := range productSlice {
		repository.Create(product)
	}

	return repository
}
