package prices

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/singleflight"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices/model"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewCoalescingController(t *testing.T) {
	demoRepo := NewDemoRepository()
	controller := NewCoalescingController(demoRepo)

	assert.NotNil(t, controller)
	assert.Equal(t, demoRepo, controller.priceRepository)
	assert.IsType(t, &singleflight.Group{}, controller.group)
}

func TestCoalescingController_GetPrices(t *testing.T) {
	t.Run("should return all prices", func(t *testing.T) {
		controller := NewCoalescingController(GenerateExampleDemoRepository())

		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/price", nil)

		controller.GetPrices(writer, request)

		res := writer.Result()
		var response []model.Price
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

		prices := GenerateExamplePriceSlice()

		if len(response) != len(prices) {
			t.Errorf("Expected count of prices is %d, got %d",
				2, len(response))
		}
	})
}

func TestCoalescingController_GetPricesByUser(t *testing.T) {
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
				priceRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest("GET", "/api/v1/price/user/abc", nil)
					request = request.WithContext(context.WithValue(request.Context(), "userId", "abc"))
					return request
				}(),
			},
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := NewCoalescingController(tt.fields.priceRepository)
			controller.GetPricesByUser(tt.args.writer, tt.args.request)
			if tt.args.writer.Code != tt.wantStatus {
				t.Errorf("Expected status code %d, got %d", tt.wantStatus, tt.args.writer.Code)
			}
		})
	}

	t.Run("Successfully get existing prices by user (expect 200 and prices)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/price/user/1", nil)
		request = request.WithContext(context.WithValue(request.Context(), "userId", "1"))

		controller := NewCoalescingController(GenerateExampleDemoRepository())

		// when
		controller.GetPricesByUser(writer, request)

		// then
		if writer.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, writer.Code)
		}

		if writer.Header().Get("Content-Type") != "application/json" {
			t.Errorf("Expected content type %s, got %s",
				"application/json", writer.Header().Get("Content-Type"))
		}

		result := writer.Result()
		var response []model.Price
		err := json.NewDecoder(result.Body).Decode(&response)
		if err != nil {
			t.Fatal(err.Error())
		}

		if len(response) != 1 {
			t.Errorf("Expected count of prices is %d, got %d",
				1, len(response))
		}

		for i, price := range response {
			if price.UserId != 1 {
				t.Errorf("Expected role of user %d, got %d", 1, response[i].UserId)
			}
		}
	})
}

func TestCoalescingController_GetPrice(t *testing.T) {
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
				priceRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer:  httptest.NewRecorder(),
				request: utils.CreatePriceRequestWithValues("GET", "/api/v1/price/abc/abc", "abc", "abc"),
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Unknown price (expect 404)",
			fields: fields{
				priceRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer:  httptest.NewRecorder(),
				request: utils.CreatePriceRequestWithValues("GET", "/api/v1/price/42/42", "42", "42"),
			},
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := NewCoalescingController(tt.fields.priceRepository)
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

		controller := NewCoalescingController(GenerateExampleDemoRepository())
		controller.GetPrice(writer, request)

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

func TestCoalescingController_PostPrice(t *testing.T) {
	type fields struct {
		priceRepository Repository
	}
	type args struct {
		writer  *httptest.ResponseRecorder
		request *http.Request
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		expectedStatus int
	}{
		{
			name: "Unauthorized (expect 401)",
			fields: fields{
				priceRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"POST",
						"/api/v1/price/4/4",
						strings.NewReader(`{"price": 0.99}`))
					ctx := context.WithValue(request.Context(), "productId", "4")
					ctx = context.WithValue(ctx, "userId", "4")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Valid create (expect 200)",
			fields: fields{
				priceRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"POST",
						"/api/v1/price/4/4",
						strings.NewReader(`{"price": 0.99}`))
					ctx := context.WithValue(request.Context(), "auth_userId", uint64(4))
					ctx = context.WithValue(ctx, "auth_userRole", int64(1))
					ctx = context.WithValue(ctx, "productId", "4")
					ctx = context.WithValue(ctx, "userId", "4")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Malformed JSON (expect 400)",
			fields: fields{
				priceRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"POST",
						"/api/v1/price/4/4",
						strings.NewReader(`{"price": 0.99`))
					ctx := context.WithValue(request.Context(), "auth_userId", uint64(4))
					ctx = context.WithValue(ctx, "auth_userRole", int64(1))
					ctx = context.WithValue(ctx, "productId", "4")
					ctx = context.WithValue(ctx, "userId", "4")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid price, incorrect Type for price (expect 400))",
			fields: fields{
				priceRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"POST",
						"/api/v1/price/4/4",
						strings.NewReader(`{"price": "0.99"}`))
					ctx := context.WithValue(request.Context(), "auth_userId", uint64(4))
					ctx = context.WithValue(ctx, "auth_userRole", int64(1))
					ctx = context.WithValue(ctx, "productId", "4")
					ctx = context.WithValue(ctx, "userId", "4")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := NewCoalescingController(tt.fields.priceRepository)
			controller.PostPrice(tt.args.writer, tt.args.request)

			if tt.args.writer.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, but got %d", tt.expectedStatus, tt.args.writer.Code)
			}
		})
	}
}

func TestCoalescingController_PutPrice(t *testing.T) {
	type fields struct {
		priceRepository Repository
	}
	type args struct {
		writer  *httptest.ResponseRecorder
		request *http.Request
	}

	tests := []struct {
		name           string
		fields         fields
		args           args
		expectedStatus int
	}{
		{
			name: "Unauthorized (expect 401)",
			fields: fields{
				priceRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/price/1/1",
						strings.NewReader(`{"price": 0.99}`))
					ctx := context.WithValue(request.Context(), "productId", "1")
					ctx = context.WithValue(ctx, "userId", "1")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Valid update (expect 200)",
			fields: fields{
				priceRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/price/1/1",
						strings.NewReader(`{"price": 10.99}`))
					ctx := context.WithValue(request.Context(), "auth_userId", uint64(1))
					ctx = context.WithValue(ctx, "auth_userRole", int64(2))
					ctx = context.WithValue(ctx, "productId", "1")
					ctx = context.WithValue(ctx, "userId", "1")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Malformed JSON (expect 400)",
			fields: fields{
				priceRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/price/2/2",
						strings.NewReader(`{"price": 6.50`))
					ctx := context.WithValue(request.Context(), "auth_userId", uint64(1))
					ctx = context.WithValue(ctx, "auth_userRole", int64(2))
					ctx = context.WithValue(ctx, "productId", "1")
					ctx = context.WithValue(ctx, "userId", "1")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Incorrect Type for Price (expect 400)",
			fields: fields{
				priceRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/price/2/2",
						strings.NewReader(`{"price": "Wrong Type"`))
					ctx := context.WithValue(request.Context(), "auth_userId", uint64(1))
					ctx = context.WithValue(ctx, "auth_userRole", int64(2))
					ctx = context.WithValue(ctx, "productId", "1")
					ctx = context.WithValue(ctx, "userId", "1")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := NewCoalescingController(tt.fields.priceRepository)
			controller.PutPrice(tt.args.writer, tt.args.request)

			if tt.args.writer.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, but got %d", tt.expectedStatus, tt.args.writer.Code)
			}
		})
	}
}

func TestCoalescingController_DeletePrice(t *testing.T) {
	type fields struct {
		priceRepository Repository
	}
	type args struct {
		writer  *httptest.ResponseRecorder
		request *http.Request
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		expectedStatus int
	}{
		{
			name: "Unauthorized (expect 401)",
			fields: fields{
				priceRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"DELETE",
						"/api/v1/price/1/1",
						nil)
					ctx := context.WithValue(request.Context(), "productId", "1")
					ctx = context.WithValue(ctx, "userId", "1")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Valid delete (expect 200)",
			fields: fields{
				priceRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"DELETE",
						"/api/v1/price/1/1",
						nil)
					ctx := context.WithValue(request.Context(), "auth_userId", uint64(1))
					ctx = context.WithValue(ctx, "auth_userRole", int64(2))
					ctx = context.WithValue(ctx, "productId", "1")
					ctx = context.WithValue(ctx, "userId", "1")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Bad non-numeric request (expect 400)",
			fields: fields{
				priceRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"DELETE",
						"/api/v1/price/abc/abc",
						nil)
					ctx := context.WithValue(request.Context(), "auth_userId", uint64(1))
					ctx = context.WithValue(ctx, "auth_userRole", int64(2))
					ctx = context.WithValue(ctx, "productId", "abc")
					ctx = context.WithValue(ctx, "userId", "abc")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Unknown product to delete (expect 500)",
			fields: fields{
				priceRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"DELETE",
						"/api/v1/price/42/42",
						nil)
					ctx := context.WithValue(request.Context(), "auth_userId", uint64(1))
					ctx = context.WithValue(ctx, "auth_userRole", int64(2))
					ctx = context.WithValue(ctx, "productId", "42")
					ctx = context.WithValue(ctx, "userId", "42")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := NewCoalescingController(tt.fields.priceRepository)
			controller.DeletePrice(tt.args.writer, tt.args.request)
			if tt.args.writer.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, tt.args.writer.Code)
			}
		})
	}
}
