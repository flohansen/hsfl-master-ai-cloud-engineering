package user

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	user "hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user/_mock"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewDefaultController(t *testing.T) {
	var mockUserRepository = user.NewMockRepository(t)
	var userRepository Repository = mockUserRepository

	userController := defaultController{
		userRepository: userRepository,
	}

	assert.NotNil(t, userController)
	assert.Equal(t, userRepository, userController.userRepository)
}

func TestDefaultController_GetUsersByRole(t *testing.T) {
	var mockUserRepository = user.NewMockRepository(t)
	var userRepository Repository = mockUserRepository

	userController := defaultController{
		userRepository: userRepository,
	}

	t.Run("Unauthorized (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/user/role/1", nil)
		request = request.WithContext(context.WithValue(request.Context(), "userRole", "1"))

		userController.GetUsersByRole(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Bad non-numeric request (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/user/role/abc", nil)
		ctx := context.WithValue(request.Context(), "auth_userId", 1)
		ctx = context.WithValue(ctx, "userRole", "abc")
		request = request.WithContext(ctx)

		userController.GetUsersByRole(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Unknown user role (expect 404)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/user/role/10", nil)
		ctx := context.WithValue(request.Context(), "auth_userId", uint64(1))
		ctx = context.WithValue(ctx, "userRole", "10")
		request = request.WithContext(ctx)

		mockUserRepository.EXPECT().FindAllByRole(model.Role(10)).Return(make([]*model.User, 0), nil)

		userController.GetUsersByRole(writer, request)

		assert.Equal(t, http.StatusNotFound, writer.Code)
	})

	t.Run("Successfully get existing users by role (expect 200 and users)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/user/role/1", nil)
		ctx := context.WithValue(request.Context(), "auth_userId", uint64(1))
		ctx = context.WithValue(ctx, "userRole", "1")
		request = request.WithContext(ctx)

		mockUserRepository.EXPECT().FindAllByRole(model.Merchant).Return([]*model.User{
			{
				Id:       1,
				Email:    "test@example.com",
				Password: []byte("12345"),
				Name:     "Test User",
				Role:     1,
			},
		}, nil)

		userController.GetUsersByRole(writer, request)

		if writer.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, writer.Code)
		}

		if writer.Header().Get("Content-Type") != "application/json" {
			t.Errorf("Expected content type %s, got %s",
				"application/json", writer.Header().Get("Content-Type"))
		}

		result := writer.Result()
		var response []model.User
		err := json.NewDecoder(result.Body).Decode(&response)
		if err != nil {
			t.Fatal(err.Error())
		}

		if len(response) != 1 {
			t.Errorf("Expected count of user is %d, got %d",
				1, len(response))
		}

		for i, user := range response {
			if user.Role != model.Merchant {
				t.Errorf("Expected role of user %d, got %d", model.Merchant, response[i].Role)
			}
		}
	})
}

func TestDefaultController_GetUser(t *testing.T) {
	var mockUserRepository = user.NewMockRepository(t)
	var userRepository Repository = mockUserRepository

	userController := defaultController{
		userRepository: userRepository,
	}

	t.Run("Unauthorized (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/user/2", nil)
		ctx := context.WithValue(request.Context(), "userId", "2")
		request = request.WithContext(ctx)

		mockUserRepository.EXPECT().FindById(uint64(2)).Return(nil, nil)

		userController.GetUser(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Valid user (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/user/1", nil)
		ctx := context.WithValue(request.Context(), "auth_userId", uint64(1))
		ctx = context.WithValue(ctx, "auth_userRole", model.Customer)
		ctx = context.WithValue(ctx, "userId", "1")
		request = request.WithContext(ctx)

		mockUserRepository.EXPECT().FindById(uint64(1)).Return(&model.User{
			Id:       1,
			Email:    "test@example.com",
			Password: []byte("12345"),
			Name:     "Test User",
			Role:     0,
		}, nil)

		userController.GetUser(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Valid user as admin (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/user/1", nil)
		ctx := context.WithValue(request.Context(), "auth_userId", uint64(2))
		ctx = context.WithValue(ctx, "auth_userRole", model.Administrator)
		ctx = context.WithValue(ctx, "userId", "1")
		request = request.WithContext(ctx)

		mockUserRepository.EXPECT().FindById(uint64(1)).Return(&model.User{
			Id:       1,
			Email:    "test@example.com",
			Password: []byte("12345"),
			Name:     "Test User",
			Role:     2,
		}, nil)

		userController.GetUser(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Bad non-numeric request (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/user/abc", nil)
		ctx := context.WithValue(request.Context(), "userId", "abc")
		request = request.WithContext(ctx)

		userController.GetUser(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Unknown user (expect 404)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/user/10", nil)
		ctx := context.WithValue(request.Context(), "auth_userId", uint64(10))
		ctx = context.WithValue(ctx, "userId", "10")
		request = request.WithContext(ctx)

		mockUserRepository.EXPECT().FindById(uint64(10)).Return(nil, errors.New(ErrorUserNotFound))

		userController.GetUser(writer, request)

		assert.Equal(t, http.StatusNotFound, writer.Code)
	})

	t.Run("Successfully get user (expect 200 and user)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/user/1", nil)
		ctx := context.WithValue(request.Context(), "auth_userId", uint64(1))
		ctx = context.WithValue(ctx, "auth_userRole", model.Customer)
		ctx = context.WithValue(ctx, "userId", "1")
		request = request.WithContext(ctx)

		mockUserRepository.EXPECT().FindById(uint64(1)).Return(&model.User{
			Id:       1,
			Email:    "test@example.com",
			Password: []byte("12345"),
			Name:     "Test User",
			Role:     0,
		}, nil)

		userController.GetUser(writer, request)

		if writer.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, writer.Code)
		}

		if writer.Header().Get("Content-Type") != "application/json" {
			t.Errorf("Expected content type %s, got %s",
				"application/json", writer.Header().Get("Content-Type"))
		}

		result := writer.Result()
		var response model.User
		err := json.NewDecoder(result.Body).Decode(&response)
		if err != nil {
			t.Fatal(err.Error())
		}

		if response.Id != 1 {
			t.Errorf("Expected id of user %d, got %d", 1, response.Id)
		}

		if response.Email != "test@example.com" {
			t.Errorf("Expected email of user %s, got %s", "ada.lovelace@gmail.com", response.Email)
		}

		if response.Name != "Test User" {
			t.Errorf("Expected name of user %s, got %s", "Ada Lovelace", response.Name)
		}

		if response.Role != model.Customer {
			t.Errorf("Got false user role")
		}
	})
}

func TestDefaultController_PutUser(t *testing.T) {
	var mockUserRepository = user.NewMockRepository(t)
	var userRepository Repository = mockUserRepository

	userController := defaultController{
		userRepository: userRepository,
	}

	t.Run("Unauthorized (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/user/1",
			strings.NewReader(`{"email": "updated@example.com", "name": "Updated user"}`))
		ctx := context.WithValue(request.Context(), "userId", "1")
		request = request.WithContext(ctx)

		userController.PutUser(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Valid user (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/user/1",
			strings.NewReader(`{"email": "updated@example.com", "name": "Updated user"}`))
		ctx := context.WithValue(request.Context(), "auth_userId", uint64(1))
		ctx = context.WithValue(ctx, "auth_userRole", model.Customer)
		ctx = context.WithValue(ctx, "userId", "1")
		request = request.WithContext(ctx)

		mockUserRepository.EXPECT().Update(&model.User{
			Id:       1,
			Email:    "updated@example.com",
			Password: nil,
			Name:     "Updated user",
			Role:     0,
		}).Return(&model.User{
			Id:       1,
			Email:    "updated@example.com",
			Password: nil,
			Name:     "Updated user",
			Role:     0,
		}, nil)

		userController.PutUser(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Valid user as admin (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/user/1",
			strings.NewReader(`{"email": "updated@example.com", "name": "Updated user"}`))
		ctx := context.WithValue(request.Context(), "auth_userId", uint64(2))
		ctx = context.WithValue(ctx, "auth_userRole", model.Administrator)
		ctx = context.WithValue(ctx, "userId", "1")
		request = request.WithContext(ctx)

		mockUserRepository.EXPECT().Update(&model.User{
			Id:       1,
			Email:    "updated@example.com",
			Password: nil,
			Name:     "Updated user",
			Role:     0,
		}).Return(&model.User{
			Id:       1,
			Email:    "updated@example.com",
			Password: nil,
			Name:     "Updated user",
			Role:     0,
		}, nil)

		userController.PutUser(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Bad non-numeric request (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/user/abc",
			strings.NewReader(`{"email": "updated@example.com", "name": "Updated user"}`))
		ctx := context.WithValue(request.Context(), "auth_userId", uint64(2))
		ctx = context.WithValue(ctx, "auth_userRole", model.Customer)
		ctx = context.WithValue(ctx, "userId", "abc")
		request = request.WithContext(ctx)

		userController.PutUser(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Unknown user (expect 404)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/user/10",
			strings.NewReader(`{"email": "updated@example.com", "name": "Updated user"}`))
		ctx := context.WithValue(request.Context(), "auth_userId", uint64(10))
		ctx = context.WithValue(ctx, "auth_userRole", model.Customer)
		ctx = context.WithValue(ctx, "userId", "10")
		request = request.WithContext(ctx)

		mockUserRepository.EXPECT().Update(&model.User{
			Id:       10,
			Email:    "updated@example.com",
			Password: nil,
			Name:     "Updated user",
			Role:     0,
		}).Return(nil, errors.New(ErrorUserNotFound))

		userController.PutUser(writer, request)

		assert.Equal(t, http.StatusInternalServerError, writer.Code)
	})
}

func TestDefaultController_DeleteUser(t *testing.T) {
	var mockUserRepository = user.NewMockRepository(t)
	var userRepository Repository = mockUserRepository

	userController := defaultController{
		userRepository: userRepository,
	}

	t.Run("Unauthorized (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/user/1", nil)
		ctx := context.WithValue(request.Context(), "userId", "1")
		request = request.WithContext(ctx)

		mockUserRepository.EXPECT().Delete(&model.User{
			Id:       1,
			Email:    "",
			Password: nil,
			Name:     "",
			Role:     0,
		}).Return(nil)

		userController.DeleteUser(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Valid user (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/user/1", nil)
		ctx := context.WithValue(request.Context(), "auth_userId", uint64(1))
		ctx = context.WithValue(ctx, "userId", "1")
		request = request.WithContext(ctx)

		mockUserRepository.EXPECT().Delete(&model.User{
			Id:       1,
			Email:    "",
			Password: nil,
			Name:     "",
			Role:     0,
		}).Return(nil)

		userController.DeleteUser(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Valid user as admin (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/user/1", nil)
		ctx := context.WithValue(request.Context(), "auth_userId", uint64(2))
		ctx = context.WithValue(ctx, "userId", "1")
		ctx = context.WithValue(ctx, "auth_userRole", model.Administrator)
		request = request.WithContext(ctx)

		mockUserRepository.EXPECT().Delete(&model.User{
			Id:       1,
			Email:    "",
			Password: nil,
			Name:     "",
			Role:     0,
		}).Return(nil)

		userController.DeleteUser(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Bad non-numeric request (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/user/abc", nil)
		ctx := context.WithValue(request.Context(), "auth_userId", uint64(2))
		ctx = context.WithValue(ctx, "userId", "abc")
		ctx = context.WithValue(ctx, "auth_userRole", model.Administrator)
		request = request.WithContext(ctx)

		userController.DeleteUser(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Unknown user (expect 404)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/user/5", nil)
		ctx := context.WithValue(request.Context(), "auth_userId", uint64(5))
		ctx = context.WithValue(ctx, "userId", "5")
		ctx = context.WithValue(ctx, "auth_userRole", model.Administrator)
		request = request.WithContext(ctx)

		mockUserRepository.EXPECT().Delete(&model.User{
			Id:       5,
			Email:    "",
			Password: nil,
			Name:     "",
			Role:     0,
		}).Return(errors.New(ErrorUserDeletion))

		userController.DeleteUser(writer, request)

		assert.Equal(t, http.StatusInternalServerError, writer.Code)
	})
}
