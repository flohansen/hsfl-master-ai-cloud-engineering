package userShoppingList

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router/middleware/auth"
	userShoppingListMock "hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList/_mock"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewDefaultController(t *testing.T) {
	var mockListRepository = userShoppingListMock.NewMockRepository(t)
	var listRepository Repository = mockListRepository

	listController := NewDefaultController(listRepository)

	assert.NotNil(t, listController)
	assert.Equal(t, listRepository, listController.userShoppingListRepository)
}

func TestDefaultController_GetLists(t *testing.T) {
	var mockListRepository = userShoppingListMock.NewMockRepository(t)
	var listRepository Repository = mockListRepository

	listController := NewDefaultController(listRepository)

	t.Run("Unauthorized (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglist/1", nil)
		ctx := context.WithValue(request.Context(), "userId", "1")
		request = request.WithContext(ctx)

		listController.GetLists(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Unauthorized, not matching userIds (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglist/1", nil)
		ctx := context.WithValue(request.Context(), "userId", "1")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		request = request.WithContext(ctx)

		listController.GetLists(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Valid request (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglist/1", nil)
		ctx := context.WithValue(request.Context(), "userId", "1")
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindAllById(uint64(1)).Return([]*model.UserShoppingList{
			{
				Id:          1,
				UserId:      1,
				Description: "Test Shopping List",
				Completed:   false,
			},
		}, nil)

		listController.GetLists(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Valid request as admin (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglist/1", nil)
		ctx := context.WithValue(request.Context(), "userId", "1")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindAllById(uint64(1)).Return([]*model.UserShoppingList{
			{
				Id:          1,
				UserId:      1,
				Description: "Test Shopping List",
				Completed:   false,
			},
		}, nil)

		listController.GetLists(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Invalid request, non-numeric userId (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglist/abc", nil)
		ctx := context.WithValue(request.Context(), "userId", "abc")
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindAllById(uint64(1)).Return([]*model.UserShoppingList{
			{
				Id:          1,
				UserId:      1,
				Description: "Test Shopping List",
				Completed:   false,
			},
		}, nil)

		listController.GetLists(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Unknown user (expect 500)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglist/10", nil)
		ctx := context.WithValue(request.Context(), "userId", "10")
		ctx = context.WithValue(ctx, "auth_userId", uint64(10))
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindAllById(uint64(10)).Return(nil, errors.New(ErrorListNotFound))

		listController.GetLists(writer, request)

		assert.Equal(t, http.StatusInternalServerError, writer.Code)
	})

	t.Run("Successfully get user shopping lists (expect 200 and shopping lists)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglist/4", nil)
		ctx := context.WithValue(request.Context(), "userId", "4")
		ctx = context.WithValue(ctx, "auth_userId", uint64(4))
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindAllById(uint64(4)).Return([]*model.UserShoppingList{
			{
				Id:          1,
				UserId:      4,
				Description: "Test Shopping List",
				Completed:   false,
			},
			{
				Id:          2,
				UserId:      4,
				Description: "Test Shopping List",
				Completed:   false,
			},
		}, nil)

		listController.GetLists(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)

		if writer.Header().Get("Content-Type") != "application/json" {
			t.Errorf("Expected content type %s, got %s",
				"application/json", writer.Header().Get("Content-Type"))
		}

		result := writer.Result()
		var response []model.UserShoppingList
		err := json.NewDecoder(result.Body).Decode(&response)
		if err != nil {
			t.Fatal(err.Error())
		}

		if len(response) != 2 {
			t.Errorf("Expected count of shopping lists is %d, got %d",
				2, len(response))
		}

		for i, list := range response {
			if list.UserId != response[i].UserId {
				t.Errorf("Expected userId of shopping lists is %d, got %d", list.UserId, response[i].UserId)
			}
		}
	})
}

func TestDefaultController_GetList(t *testing.T) {
	var mockListRepository = userShoppingListMock.NewMockRepository(t)
	var listRepository Repository = mockListRepository

	listController := NewDefaultController(listRepository)

	t.Run("Unauthorized (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglist/1/2", nil)
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "userId", "2")
		request = request.WithContext(ctx)

		listController.GetList(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Unauthorized, not matching userIds (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglist/1/2", nil)
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "userId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		request = request.WithContext(ctx)

		listController.GetList(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Valid request (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglist/1/2", nil)
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "userId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindByIds(uint64(2), uint64(1)).Return(&model.UserShoppingList{
			Id:          1,
			UserId:      2,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		listController.GetList(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Valid request as admin (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglist/1/2", nil)
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "userId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindByIds(uint64(2), uint64(1)).Return(&model.UserShoppingList{
			Id:          1,
			UserId:      2,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		listController.GetList(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Invalid request, non-numeric userId (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglist/abc/abc", nil)
		ctx := context.WithValue(request.Context(), "listId", "abc")
		ctx = context.WithValue(ctx, "userId", "abc")
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		listController.GetList(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Unknown shopping list (expect 404)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglist/10/10", nil)
		ctx := context.WithValue(request.Context(), "listId", "10")
		ctx = context.WithValue(ctx, "userId", "10")
		ctx = context.WithValue(ctx, "auth_userId", uint64(10))
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindByIds(uint64(10), uint64(10)).Return(nil, errors.New(ErrorListNotFound))

		listController.GetList(writer, request)

		assert.Equal(t, http.StatusNotFound, writer.Code)
	})
}

func TestDefaultController_PutList(t *testing.T) {
	var mockListRepository = userShoppingListMock.NewMockRepository(t)
	var listRepository Repository = mockListRepository

	listController := NewDefaultController(listRepository)

	t.Run("Unauthorized (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/shoppinglist/1/2",
			strings.NewReader(`{"description": "Updated list"}`))
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "userId", "2")
		request = request.WithContext(ctx)

		listController.PutList(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Unauthorized, not matching userIds (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/shoppinglist/1/2",
			strings.NewReader(`{"description": "Updated list"}`))
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "userId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		request = request.WithContext(ctx)

		listController.PutList(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Valid request (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/shoppinglist/1/2",
			strings.NewReader(`{"description": "Updated list"}`))
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "userId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().Update(&model.UserShoppingList{
			Id:          1,
			UserId:      2,
			Description: "Updated list",
			Completed:   false,
		}).Return(&model.UserShoppingList{
			Id:          1,
			UserId:      2,
			Description: "Updated list",
			Completed:   false,
		}, nil)

		listController.PutList(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Valid request as admin (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/shoppinglist/1/2",
			strings.NewReader(`{"description": "Updated list"}`))
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "userId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().Update(&model.UserShoppingList{
			Id:          1,
			UserId:      2,
			Description: "Updated list",
			Completed:   false,
		}).Return(&model.UserShoppingList{
			Id:          1,
			UserId:      2,
			Description: "Updated list",
			Completed:   false,
		}, nil)

		listController.PutList(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Invalid request, non-numeric userId (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/shoppinglist/abc/abc",
			strings.NewReader(`{"description": "Updated list"}`))
		ctx := context.WithValue(request.Context(), "listId", "abc")
		ctx = context.WithValue(ctx, "userId", "abc")
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		request = request.WithContext(ctx)

		listController.PutList(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Invalid request, non-numeric userId (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/shoppinglist/abc/abc",
			strings.NewReader(`{"description": "Updated list"}`))
		ctx := context.WithValue(request.Context(), "listId", "abc")
		ctx = context.WithValue(ctx, "userId", "abc")
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		request = request.WithContext(ctx)

		listController.PutList(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Malformed JSON (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/shoppinglist/abc/abc",
			strings.NewReader(`{"description": "Updated list"`))
		ctx := context.WithValue(request.Context(), "listId", "abc")
		ctx = context.WithValue(ctx, "userId", "abc")
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		request = request.WithContext(ctx)

		listController.PutList(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Unknown shopping list (expect 500)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/shoppinglist/10/1",
			strings.NewReader(`{"description": "Updated list", "completed": false}`))
		ctx := context.WithValue(request.Context(), "listId", "10")
		ctx = context.WithValue(ctx, "userId", "1")
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().Update(&model.UserShoppingList{
			Id:          10,
			UserId:      1,
			Description: "Updated list",
			Completed:   false,
		}).Return(nil, errors.New(ErrorListNotFound))

		listController.PutList(writer, request)

		assert.Equal(t, http.StatusInternalServerError, writer.Code)
	})
}

func TestDefaultController_PostList(t *testing.T) {
	var mockListRepository = userShoppingListMock.NewMockRepository(t)
	var listRepository Repository = mockListRepository

	listController := NewDefaultController(listRepository)

	t.Run("Unauthorized (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/v1/shoppinglist/2",
			strings.NewReader(`{"description": "New list", "checked": false}`))
		ctx := context.WithValue(request.Context(), "userId", "2")
		request = request.WithContext(ctx)

		listController.PostList(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Unauthorized, not matching userIds (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/v1/shoppinglist/2",
			strings.NewReader(`{"description": "New list", "checked": false}`))
		ctx := context.WithValue(request.Context(), "userId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		request = request.WithContext(ctx)

		listController.PostList(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Valid request (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/v1/shoppinglist/2",
			strings.NewReader(`{"description": "New list", "checked": false}`))
		ctx := context.WithValue(request.Context(), "userId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().Create(&model.UserShoppingList{
			UserId:      2,
			Description: "New list",
			Completed:   false,
		}).Return(&model.UserShoppingList{
			Id:          1,
			UserId:      2,
			Description: "New list",
			Completed:   false,
		}, nil)

		listController.PostList(writer, request)

		assert.Equal(t, http.StatusCreated, writer.Code)
	})

	t.Run("Valid request as admin(expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/v1/shoppinglist/2",
			strings.NewReader(`{"description": "New list", "checked": false}`))
		ctx := context.WithValue(request.Context(), "userId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().Create(&model.UserShoppingList{
			UserId:      2,
			Description: "New list",
			Completed:   false,
		}).Return(&model.UserShoppingList{
			Id:          1,
			UserId:      2,
			Description: "New list",
			Completed:   false,
		}, nil)

		listController.PostList(writer, request)

		assert.Equal(t, http.StatusCreated, writer.Code)
	})

	t.Run("Invalid request, non-numeric userId (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/v1/shoppinglist/abc",
			strings.NewReader(`{"description": "New list", "checked": false}`))
		ctx := context.WithValue(request.Context(), "userId", "abc")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		request = request.WithContext(ctx)

		listController.PostList(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Malformed JSON (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/v1/shoppinglist/2",
			strings.NewReader(`{"description": "New list", "checked": false`))
		ctx := context.WithValue(request.Context(), "userId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		request = request.WithContext(ctx)

		listController.PostList(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})
}

func TestDefaultController_Delete(t *testing.T) {
	var mockListRepository = userShoppingListMock.NewMockRepository(t)
	var listRepository Repository = mockListRepository

	listController := NewDefaultController(listRepository)

	t.Run("Unauthorized (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/shoppinglist/1/2", nil)
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "userId", "2")
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&model.UserShoppingList{
			Id:          1,
			UserId:      2,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		listController.DeleteList(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Unauthorized, not matching userIds (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/shoppinglist/1/2", nil)
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "userId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&model.UserShoppingList{
			Id:          1,
			UserId:      2,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		listController.DeleteList(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Valid request (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/shoppinglist/1/2", nil)
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "userId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&model.UserShoppingList{
			Id:          1,
			UserId:      2,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		mockListRepository.EXPECT().Delete(&model.UserShoppingList{Id: 1}).Return(nil)

		listController.DeleteList(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Valid request as admin (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/shoppinglist/1/2", nil)
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "userId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&model.UserShoppingList{
			Id:          1,
			UserId:      2,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		mockListRepository.EXPECT().Delete(&model.UserShoppingList{Id: 1}).Return(nil)

		listController.DeleteList(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Invalid request, non-numeric listId (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/shoppinglist/abc/abc", nil)
		ctx := context.WithValue(request.Context(), "listId", "abc")
		ctx = context.WithValue(ctx, "userId", "abc")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		request = request.WithContext(ctx)

		listController.DeleteList(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Unknown shopping list (expect 404)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/shoppinglist/10/2", nil)
		ctx := context.WithValue(request.Context(), "listId", "10")
		ctx = context.WithValue(ctx, "userId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(10)).Return(nil, errors.New(ErrorListNotFound))

		listController.DeleteList(writer, request)

		assert.Equal(t, http.StatusNotFound, writer.Code)
	})
}
