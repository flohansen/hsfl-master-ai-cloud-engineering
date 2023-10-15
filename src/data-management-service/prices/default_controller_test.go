package prices

import (
	"context"
	"encoding/json"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/data-management-service/prices/model"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewDefaultController(t *testing.T) {
	type args struct {
		priceRepository Repository
	}
	tests := []struct {
		name string
		args args
		want *defaultController
	}{
		{
			name: "Test construction with DemoRepository",
			args: args{priceRepository: NewDemoRepository()},
			want: &defaultController{priceRepository: NewDemoRepository()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDefaultController(tt.args.priceRepository); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDefaultController() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultController_DeletePrice(t *testing.T) {
	type fields struct {
		priceRepository Repository
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
			name: "Successfully delete existing price (expect 200)",
			fields: fields{
				priceRepository: setupMockRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest("DELETE", "/api/v1/price/1/1", nil)
					request = request.WithContext(context.WithValue(request.Context(), "productId", "1"))
					request = request.WithContext(context.WithValue(request.Context(), "userId", "1"))
					return request
				}(),
			},

			wantStatus: http.StatusOK,
		},
		{
			name: "Bad non-numeric request (expect 400)",
			fields: fields{
				priceRepository: setupMockRepository(),
			},
			args: args{
				writer:  httptest.NewRecorder(),
				request: createRequestWithValues("DELETE", "/api/v1/price/abc/abc", "abc", "abc"),
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Unknown product to delete (expect 500)",
			fields: fields{
				priceRepository: setupMockRepository(),
			},
			args: args{
				writer:  httptest.NewRecorder(),
				request: createRequestWithValues("DELETE", "/api/v1/price/42/42", "42", "42"),
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := defaultController{
				priceRepository: tt.fields.priceRepository,
			}
			controller.DeletePrice(tt.args.writer, tt.args.request)
			if tt.args.writer.Code != tt.wantStatus {
				t.Errorf("Expected status code %d, got %d", tt.wantStatus, tt.args.writer.Code)
			}
		})
	}
}

func TestDefaultController_GetPrice(t *testing.T) {
	type fields struct {
		priceRepository Repository
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
				priceRepository: setupMockRepository(),
			},
			args: args{
				writer:  httptest.NewRecorder(),
				request: createRequestWithValues("GET", "/api/v1/price/abc/abc", "abc", "abc"),
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Unknown price (expect 404)",
			fields: fields{
				priceRepository: setupMockRepository(),
			},
			args: args{
				writer:  httptest.NewRecorder(),
				request: createRequestWithValues("GET", "/api/v1/price/42/42", "42", "42"),
			},
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := defaultController{
				priceRepository: tt.fields.priceRepository,
			}
			controller.GetPrice(tt.args.writer, tt.args.request)
			if tt.args.writer.Code != tt.wantStatus {
				t.Errorf("Expected status code %d, got %d", tt.wantStatus, tt.args.writer.Code)
			}
		})
	}

	t.Run("Successfully get existing price (expect 200 and price)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/price/1/1", nil)
		request = request.WithContext(context.WithValue(request.Context(), "productId", "1"))
		request = request.WithContext(context.WithValue(request.Context(), "userId", "1"))

		controller := defaultController{
			priceRepository: setupMockRepository(),
		}

		// when
		controller.GetPrice(writer, request)

		// then
		if writer.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, writer.Code)
		}

		if writer.Header().Get("Content-Type") != "application/json" {
			t.Errorf("Expected content type %s, got %s",
				"application/json", writer.Header().Get("Content-Type"))
		}

		result := writer.Result()
		var response model.Price
		err := json.NewDecoder(result.Body).Decode(&response)
		if err != nil {
			t.Fatal(err.Error())
		}

		if response.ProductId != 1 {
			t.Errorf("Expected product id of price %d, got %d", 1, response.ProductId)
		}

		if response.UserId != 1 {
			t.Errorf("Expected user id of product %d, got %d", 1, response.UserId)
		}

		if response.Price != 2.99 {
			t.Errorf("Expected ean of product %f, got %f", 2.99, response.Price)
		}

	})
}

func createRequestWithValues(method, path string, productId, userId string) *http.Request {
	request := httptest.NewRequest(method, path, nil)
	ctx := request.Context()
	ctx = context.WithValue(ctx, "productId", productId)
	ctx = context.WithValue(ctx, "userId", userId)
	request = request.WithContext(ctx)
	return request
}

func setupMockRepository() Repository {
	repository := NewDemoRepository()
	pricesSlice := []*model.Price{
		{
			UserId:    1,
			ProductId: 1,
			Price:     2.99,
		},
		{
			UserId:    2,
			ProductId: 2,
			Price:     5.99,
		},
	}
	for _, price := range pricesSlice {
		repository.Create(price)
	}

	return repository
}
