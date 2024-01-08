package prices

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/singleflight"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router/middleware/auth"
	prices "hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices/_mock"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewCoalescingController(t *testing.T) {
	var mockPriceRepository = prices.NewMockRepository(t)
	var priceRepository Repository = mockPriceRepository

	priceController := NewCoalescingController(priceRepository)

	assert.NotNil(t, priceController)
	assert.Equal(t, priceRepository, priceController.priceRepository)
	assert.IsType(t, &singleflight.Group{}, priceController.group)
}

func TestCoalescingController_GetPrices(t *testing.T) {
	var mockPriceRepository = prices.NewMockRepository(t)
	var priceRepository Repository = mockPriceRepository

	priceController := NewCoalescingController(priceRepository)

	t.Run("Successfully return all prices (expect 200 and prices)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/price", nil)

		mockPriceRepository.EXPECT().FindAll().Return([]*model.Price{
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
		}, nil)

		priceController.GetPrices(writer, request)

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

		if len(response) != 2 {
			t.Errorf("Expected count of prices is %d, got %d",
				2, len(response))
		}
	})
}

func TestCoalescingController_GetPricesByUser(t *testing.T) {
	var mockPriceRepository = prices.NewMockRepository(t)
	var priceRepository Repository = mockPriceRepository

	priceController := NewCoalescingController(priceRepository)

	t.Run("Bad non-numeric request (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/price/user/abc", nil)
		ctx := context.WithValue(request.Context(), "userId", "abc")
		request = request.WithContext(ctx)

		priceController.GetPricesByUser(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Unknown userId (expect 404)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/price/user/10", nil)
		ctx := context.WithValue(request.Context(), "userId", "10")
		request = request.WithContext(ctx)

		mockPriceRepository.EXPECT().FindAllByUser(uint64(10)).Return(nil, errors.New(ErrorPriceNotFound))

		priceController.GetPricesByUser(writer, request)

		assert.Equal(t, http.StatusNotFound, writer.Code)
	})

	t.Run("Successfully get existing prices by user (expect 200 and prices)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/price/user/1", nil)
		request = request.WithContext(context.WithValue(request.Context(), "userId", "1"))

		mockPriceRepository.EXPECT().FindAllByUser(uint64(1)).Return([]*model.Price{
			{
				UserId:    1,
				ProductId: 1,
				Price:     2.99,
			},
		}, nil)

		priceController.GetPricesByUser(writer, request)

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
	var mockPriceRepository = prices.NewMockRepository(t)
	var priceRepository Repository = mockPriceRepository

	priceController := NewCoalescingController(priceRepository)

	t.Run("Bad non-numeric request (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/price/abc/abc", nil)
		ctx := context.WithValue(request.Context(), "userId", "abc")
		ctx = context.WithValue(ctx, "productId", "abc")
		request = request.WithContext(ctx)

		priceController.GetPrice(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Unknown price (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/price/10/10", nil)
		ctx := context.WithValue(request.Context(), "userId", "10")
		ctx = context.WithValue(ctx, "productId", "10")
		request = request.WithContext(ctx)

		mockPriceRepository.EXPECT().FindByIds(uint64(10), uint64(10)).Return(nil, errors.New(ErrorPriceNotFound))

		priceController.GetPrice(writer, request)

		assert.Equal(t, http.StatusNotFound, writer.Code)
	})

	t.Run("Successfully get existing price (expect 200 and price)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/price/1/1", nil)
		request = request.WithContext(context.WithValue(request.Context(), "productId", "1"))
		request = request.WithContext(context.WithValue(request.Context(), "userId", "1"))

		mockPriceRepository.EXPECT().FindByIds(uint64(1), uint64(1)).Return(&model.Price{
			UserId:    1,
			ProductId: 1,
			Price:     2.99,
		}, nil)

		priceController.GetPrice(writer, request)

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
	var mockPriceRepository = prices.NewMockRepository(t)
	var priceRepository Repository = mockPriceRepository

	priceController := NewCoalescingController(priceRepository)

	t.Run("Unauthorized (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/v1/price/4/4",
			strings.NewReader(`{"price": 0.99}`))
		ctx := context.WithValue(request.Context(), "userId", "4")
		ctx = context.WithValue(ctx, "productId", "4")
		request = request.WithContext(ctx)

		priceController.PostPrice(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Valid create (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/v1/price/4/4",
			strings.NewReader(`{"price": 0.99}`))
		ctx := context.WithValue(request.Context(), "userId", "4")
		ctx = context.WithValue(ctx, "productId", "4")
		ctx = context.WithValue(ctx, "auth_userRole", auth.Merchant)
		ctx = context.WithValue(ctx, "auth_userId", uint64(4))
		request = request.WithContext(ctx)

		mockPriceRepository.EXPECT().Create(&model.Price{
			UserId:    4,
			ProductId: 4,
			Price:     0.99,
		}).Return(&model.Price{
			UserId:    4,
			ProductId: 4,
			Price:     0.99,
		}, nil)

		priceController.PostPrice(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Valid create as admin (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/v1/price/4/4",
			strings.NewReader(`{"price": 0.99}`))
		ctx := context.WithValue(request.Context(), "userId", "4")
		ctx = context.WithValue(ctx, "productId", "4")
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		request = request.WithContext(ctx)

		mockPriceRepository.EXPECT().Create(&model.Price{
			UserId:    4,
			ProductId: 4,
			Price:     0.99,
		}).Return(&model.Price{
			UserId:    4,
			ProductId: 4,
			Price:     0.99,
		}, nil)

		priceController.PostPrice(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Malformed JSON (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/v1/price/4/4",
			strings.NewReader(`{"price": 0.99`))
		ctx := context.WithValue(request.Context(), "userId", "4")
		ctx = context.WithValue(ctx, "productId", "4")
		ctx = context.WithValue(ctx, "auth_userRole", auth.Merchant)
		ctx = context.WithValue(ctx, "auth_userId", uint64(4))
		request = request.WithContext(ctx)

		priceController.PostPrice(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Invalid price (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/v1/price/4/4",
			strings.NewReader(`{"price": "0.99"}`))
		ctx := context.WithValue(request.Context(), "userId", "4")
		ctx = context.WithValue(ctx, "productId", "4")
		ctx = context.WithValue(ctx, "auth_userRole", auth.Merchant)
		ctx = context.WithValue(ctx, "auth_userId", uint64(4))
		request = request.WithContext(ctx)

		priceController.PostPrice(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})
}

func TestCoalescingController_PutPrice(t *testing.T) {
	var mockPriceRepository = prices.NewMockRepository(t)
	var priceRepository Repository = mockPriceRepository

	priceController := NewCoalescingController(priceRepository)

	t.Run("Unauthorized (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/price/1/1",
			strings.NewReader(`{"price": 10.99}`))
		ctx := context.WithValue(request.Context(), "userId", "1")
		ctx = context.WithValue(ctx, "productId", "1")
		request = request.WithContext(ctx)

		priceController.PutPrice(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Valid update (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/price/1/1",
			strings.NewReader(`{"price": 10.99}`))
		ctx := context.WithValue(request.Context(), "auth_userId", uint64(1))
		ctx = context.WithValue(ctx, "auth_userRole", int64(2))
		ctx = context.WithValue(ctx, "productId", "1")
		ctx = context.WithValue(ctx, "userId", "1")
		request = request.WithContext(ctx)

		mockPriceRepository.EXPECT().Update(&model.Price{
			UserId:    1,
			ProductId: 1,
			Price:     10.99,
		}).Return(&model.Price{
			UserId:    1,
			ProductId: 1,
			Price:     10.99,
		}, nil)

		priceController.PutPrice(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Valid update as admin (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/price/1/1",
			strings.NewReader(`{"price": 10.99}`))
		ctx := context.WithValue(request.Context(), "auth_userId", uint64(10))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		ctx = context.WithValue(ctx, "productId", "1")
		ctx = context.WithValue(ctx, "userId", "1")
		request = request.WithContext(ctx)

		mockPriceRepository.EXPECT().Update(&model.Price{
			UserId:    1,
			ProductId: 1,
			Price:     10.99,
		}).Return(&model.Price{
			UserId:    1,
			ProductId: 1,
			Price:     10.99,
		}, nil)

		priceController.PutPrice(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Malformed JSON (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/price/1/1",
			strings.NewReader(`{"price": 10.99`))
		ctx := context.WithValue(request.Context(), "auth_userId", uint64(1))
		ctx = context.WithValue(ctx, "auth_userRole", int64(2))
		ctx = context.WithValue(ctx, "productId", "1")
		ctx = context.WithValue(ctx, "userId", "1")
		request = request.WithContext(ctx)

		priceController.PutPrice(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Incorrect Type for Price (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/price/1/1",
			strings.NewReader(`{"price": "10.99"}`))
		ctx := context.WithValue(request.Context(), "auth_userId", uint64(1))
		ctx = context.WithValue(ctx, "auth_userRole", int64(2))
		ctx = context.WithValue(ctx, "productId", "1")
		ctx = context.WithValue(ctx, "userId", "1")
		request = request.WithContext(ctx)

		priceController.PutPrice(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Unknown price (expect 500)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/price/10/10",
			strings.NewReader(`{"price": 10.99}`))
		ctx := context.WithValue(request.Context(), "auth_userId", uint64(10))
		ctx = context.WithValue(ctx, "auth_userRole", int64(2))
		ctx = context.WithValue(ctx, "productId", "10")
		ctx = context.WithValue(ctx, "userId", "10")
		request = request.WithContext(ctx)

		mockPriceRepository.EXPECT().Update(&model.Price{
			UserId:    10,
			ProductId: 10,
			Price:     10.99,
		}).Return(nil, errors.New(ErrorPriceNotFound))

		priceController.PutPrice(writer, request)

		assert.Equal(t, http.StatusInternalServerError, writer.Code)
	})
}

func TestCoalescingController_DeletePrice(t *testing.T) {
	var mockPriceRepository = prices.NewMockRepository(t)
	var priceRepository Repository = mockPriceRepository

	priceController := NewCoalescingController(priceRepository)

	t.Run("Unauthorized (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/price/1/1", nil)
		ctx := context.WithValue(request.Context(), "userId", "1")
		ctx = context.WithValue(ctx, "productId", "1")
		request = request.WithContext(ctx)

		priceController.DeletePrice(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Valid delete (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/price/1/1", nil)
		ctx := context.WithValue(request.Context(), "userId", "1")
		ctx = context.WithValue(ctx, "productId", "1")
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Merchant)
		request = request.WithContext(ctx)

		mockPriceRepository.EXPECT().Delete(&model.Price{
			UserId:    1,
			ProductId: 1,
		}).Return(nil)

		priceController.DeletePrice(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Valid delete as admin (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/price/1/1", nil)
		ctx := context.WithValue(request.Context(), "userId", "1")
		ctx = context.WithValue(ctx, "productId", "1")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		mockPriceRepository.EXPECT().Delete(&model.Price{
			UserId:    1,
			ProductId: 1,
		}).Return(nil)

		priceController.DeletePrice(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Bad non-numeric request (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/price/abc/abc", nil)
		ctx := context.WithValue(request.Context(), "userId", "abc")
		ctx = context.WithValue(ctx, "productId", "abc")
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Merchant)
		request = request.WithContext(ctx)

		priceController.DeletePrice(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Unknown price (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/price/10/10", nil)
		ctx := context.WithValue(request.Context(), "userId", "10")
		ctx = context.WithValue(ctx, "productId", "10")
		ctx = context.WithValue(ctx, "auth_userId", uint64(10))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Merchant)
		request = request.WithContext(ctx)

		mockPriceRepository.EXPECT().Delete(&model.Price{
			UserId:    10,
			ProductId: 10,
		}).Return(errors.New(ErrorPriceNotFound))

		priceController.DeletePrice(writer, request)

		assert.Equal(t, http.StatusInternalServerError, writer.Code)
	})
}
