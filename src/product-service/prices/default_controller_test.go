package prices

import (
	"context"
	"encoding/json"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices/model"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices/utils"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
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
				priceRepository: GenerateExampleDemoRepository(),
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
				priceRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer:  httptest.NewRecorder(),
				request: utils.CreatePriceRequestWithValues("DELETE", "/api/v1/price/abc/abc", "abc", "abc"),
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Unknown product to delete (expect 500)",
			fields: fields{
				priceRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer:  httptest.NewRecorder(),
				request: utils.CreatePriceRequestWithValues("DELETE", "/api/v1/price/42/42", "42", "42"),
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

func TestDefaultController_GetPrices(t *testing.T) {
	t.Run("should return all prices", func(t *testing.T) {
		controller := defaultController{
			priceRepository: GenerateExampleDemoRepository(),
		}

		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/price", nil)

		// Test request
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

func TestDefaultController_GetPricesByUser(t *testing.T) {
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
			controller := defaultController{
				priceRepository: tt.fields.priceRepository,
			}
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

		controller := defaultController{
			priceRepository: GenerateExampleDemoRepository(),
		}

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
			priceRepository: GenerateExampleDemoRepository(),
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

func TestDefaultController_PostPrice(t *testing.T) {
	type fields struct {
		priceRepository Repository
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
			name: "Valid Price",
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
					request = request.WithContext(context.WithValue(request.Context(), "productId", "4"))
					request = request.WithContext(context.WithValue(request.Context(), "userId", "4"))
					return request
				}(),
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "",
		},
		{
			name: "Malformed JSON",
			fields: fields{
				priceRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"POST",
						"/api/v1/price/5/5",
						strings.NewReader(`{"price": 0.99`))
					request = request.WithContext(context.WithValue(request.Context(), "productId", "5"))
					request = request.WithContext(context.WithValue(request.Context(), "userId", "5"))
					return request
				}(),
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "",
		},
		{
			name: "Invalid price, incorrect Type for price (Non-numeric)",
			fields: fields{
				priceRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"POST",
						"/api/v1/price/5/5",
						strings.NewReader(`{"price": "0.99"}`))
					request = request.WithContext(context.WithValue(request.Context(), "productId", "5"))
					request = request.WithContext(context.WithValue(request.Context(), "userId", "5"))
					return request
				}(),
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := defaultController{
				priceRepository: tt.fields.priceRepository,
			}
			controller.PostPrice(tt.args.writer, tt.args.request)

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

func TestDefaultController_PutPrice(t *testing.T) {
	type fields struct {
		priceRepository Repository
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
				priceRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/price/1/1",
						strings.NewReader(`{"userId": 1, "productId": 1, "price": 10.99}`))
					request = request.WithContext(context.WithValue(request.Context(), "productId", "1"))
					request = request.WithContext(context.WithValue(request.Context(), "userId", "1"))
					return request
				}(),
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "",
		},
		{
			name: "Valid Update (Partly Fields)",
			fields: fields{
				priceRepository: GenerateExampleDemoRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/price/2/2",
						strings.NewReader(`{"price": 6.50}`))
					request = request.WithContext(context.WithValue(request.Context(), "productId", "2"))
					request = request.WithContext(context.WithValue(request.Context(), "userId", "2"))
					return request
				}(),
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "",
		},
		{
			name: "Malformed JSON",
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
					request = request.WithContext(context.WithValue(request.Context(), "productId", "2"))
					request = request.WithContext(context.WithValue(request.Context(), "userId", "2"))
					return request
				}(),
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "",
		},
		{
			name:   "Incorrect Type for Price (Non-numeric)",
			fields: fields{
				// Set up your repository mock or test double here if needed
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/price/2/2",
						strings.NewReader(`{"price": "Wrong Type"`))
					request = request.WithContext(context.WithValue(request.Context(), "productId", "2"))
					request = request.WithContext(context.WithValue(request.Context(), "userId", "2"))
					return request
				}(),
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := defaultController{
				priceRepository: tt.fields.priceRepository,
			}
			controller.PutPrice(tt.args.writer, tt.args.request)

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
