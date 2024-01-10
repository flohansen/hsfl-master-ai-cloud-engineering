package userShoppingListEntry

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router/middleware/auth"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList"
	userShoppingListMock "hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList/_mock"
	listModel "hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList/model"
	userShoppingListEntryMock "hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingListEntry/_mock"
	entryModel "hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingListEntry/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewDefaultController(t *testing.T) {
	var mockEntryRepository = userShoppingListEntryMock.NewMockRepository(t)
	var entryRepository Repository = mockEntryRepository

	var mockListRepository = userShoppingListMock.NewMockRepository(t)
	var listRepository userShoppingList.Repository = mockListRepository

	entryController := NewDefaultController(entryRepository, listRepository)

	assert.NotNil(t, entryController)
	assert.Equal(t, entryRepository, entryController.userShoppingListEntryRepository)
	assert.Equal(t, listRepository, entryController.userShoppingListRepository)
}

func TestDefaultController_GetEntries(t *testing.T) {
	var mockEntryRepository = userShoppingListEntryMock.NewMockRepository(t)
	var entryRepository Repository = mockEntryRepository

	var mockListRepository = userShoppingListMock.NewMockRepository(t)
	var listRepository userShoppingList.Repository = mockListRepository

	entryController := NewDefaultController(entryRepository, listRepository)

	t.Run("Unauthorized (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglistentries/1", nil)
		ctx := context.WithValue(request.Context(), "listId", "1")
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&listModel.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		entryController.GetEntries(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Unauthorized not matching userIds (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglistentries/1", nil)
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&listModel.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		entryController.GetEntries(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Valid request (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglistentries/1", nil)
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&listModel.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		mockEntryRepository.EXPECT().FindAll(uint64(1)).Return([]*entryModel.UserShoppingListEntry{
			{
				ShoppingListId: 1,
				ProductId:      1,
				Count:          1,
				Note:           "",
				Checked:        false,
			},
		}, nil)

		entryController.GetEntries(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Valid request as admin (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglistentries/1", nil)
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&listModel.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		mockEntryRepository.EXPECT().FindAll(uint64(1)).Return([]*entryModel.UserShoppingListEntry{
			{
				ShoppingListId: 1,
				ProductId:      1,
				Count:          1,
				Note:           "",
				Checked:        false,
			},
		}, nil)

		entryController.GetEntries(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Invalid request, non-numeric listId (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglistentries/abc", nil)
		ctx := context.WithValue(request.Context(), "listId", "abc")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		entryController.GetEntries(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Unknown list (expect 404)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglistentries/10", nil)
		ctx := context.WithValue(request.Context(), "listId", "10")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(10)).Return(nil, errors.New(userShoppingList.ErrorListNotFound))

		entryController.GetEntries(writer, request)

		assert.Equal(t, http.StatusNotFound, writer.Code)
	})

	t.Run("Successfully get existing entries (expect 200 and entries)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglistentries/2", nil)
		ctx := context.WithValue(request.Context(), "listId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(2)).Return(&listModel.UserShoppingList{
			Id:          2,
			UserId:      1,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		mockEntryRepository.EXPECT().FindAll(uint64(2)).Return([]*entryModel.UserShoppingListEntry{
			{
				ShoppingListId: 2,
				ProductId:      1,
				Count:          1,
				Note:           "",
				Checked:        false,
			},
			{
				ShoppingListId: 2,
				ProductId:      2,
				Count:          10,
				Note:           "",
				Checked:        false,
			},
			{
				ShoppingListId: 2,
				ProductId:      3,
				Count:          10,
				Note:           "",
				Checked:        false,
			},
		}, nil)

		entryController.GetEntries(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)

		if writer.Header().Get("Content-Type") != "application/json" {
			t.Errorf("Expected content type %s, got %s",
				"application/json", writer.Header().Get("Content-Type"))
		}

		res := writer.Result()
		var response []entryModel.UserShoppingListEntry
		err := json.NewDecoder(res.Body).Decode(&response)

		if err != nil {
			t.Error(err)
		}

		if len(response) != 3 {
			t.Errorf("Expected count of shopping list entries is %d, got %d",
				3, len(response))
		}

		for i, entry := range response {
			if entry.ShoppingListId != 2 {
				t.Errorf("Expected shopping list id of entry %d, got %d", 2, response[i].ShoppingListId)
			}
		}
	})
}

func TestDefaultController_GetEntry(t *testing.T) {
	var mockEntryRepository = userShoppingListEntryMock.NewMockRepository(t)
	var entryRepository Repository = mockEntryRepository

	var mockListRepository = userShoppingListMock.NewMockRepository(t)
	var listRepository userShoppingList.Repository = mockListRepository

	entryController := NewDefaultController(entryRepository, listRepository)

	t.Run("Unauthorized (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglistentries/1/2", nil)
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "productId", "2")
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&listModel.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		mockEntryRepository.EXPECT().FindByIds(uint64(1), uint64(2)).Return(nil, nil)

		entryController.GetEntry(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Unauthorized, not matching userIds (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglistentries/1/2", nil)
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "productId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&listModel.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		mockEntryRepository.EXPECT().FindByIds(uint64(1), uint64(2)).Return(nil, nil)

		entryController.GetEntry(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Valid request (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglistentries/1/2", nil)
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "productId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&listModel.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		mockEntryRepository.EXPECT().FindByIds(uint64(1), uint64(2)).Return(&entryModel.UserShoppingListEntry{
			ShoppingListId: 1,
			ProductId:      2,
			Count:          10,
			Note:           "",
			Checked:        false,
		}, nil)

		entryController.GetEntry(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Valid request as admin (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglistentries/1/2", nil)
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "productId", "2")
		ctx = context.WithValue(ctx, "auth_userId", int64(2))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&listModel.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		mockEntryRepository.EXPECT().FindByIds(uint64(1), uint64(2)).Return(&entryModel.UserShoppingListEntry{
			ShoppingListId: 1,
			ProductId:      2,
			Count:          10,
			Note:           "",
			Checked:        false,
		}, nil)

		entryController.GetEntry(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Invalid request, non-numeric listId and productId (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglistentries/abc/abc", nil)
		ctx := context.WithValue(request.Context(), "listId", "abc")
		ctx = context.WithValue(ctx, "productId", "abc")
		ctx = context.WithValue(ctx, "auth_userId", int64(2))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		entryController.GetEntry(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Unknown list (expect 404)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglistentries/10/1", nil)
		ctx := context.WithValue(request.Context(), "listId", "10")
		ctx = context.WithValue(ctx, "productId", "1")
		ctx = context.WithValue(ctx, "auth_userId", int64(2))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(10)).Return(nil, errors.New(userShoppingList.ErrorListNotFound))

		entryController.GetEntry(writer, request)

		assert.Equal(t, http.StatusNotFound, writer.Code)
	})

	t.Run("Unknown entry (expect 500)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglistentries/1/10", nil)
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "productId", "10")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&listModel.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		mockEntryRepository.EXPECT().FindByIds(uint64(1), uint64(10)).Return(nil, errors.New(ErrorEntryNotFound))

		entryController.GetEntry(writer, request)

		assert.Equal(t, http.StatusNotFound, writer.Code)
	})
}

func TestDefaultController_PostEntry(t *testing.T) {
	var mockEntryRepository = userShoppingListEntryMock.NewMockRepository(t)
	var entryRepository Repository = mockEntryRepository

	var mockListRepository = userShoppingListMock.NewMockRepository(t)
	var listRepository userShoppingList.Repository = mockListRepository

	entryController := NewDefaultController(entryRepository, listRepository)

	t.Run("Unauthorized (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/v1/shoppinglistentries/1/2",
			strings.NewReader(`{"count": 2, "note": "Test entry", "checked": false}`))
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "productId", "2")
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&listModel.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		entryController.PostEntry(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Unauthorized, not matching userId (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/v1/shoppinglistentries/1/2",
			strings.NewReader(`{"count": 2, "note": "Test entry", "checked": false}`))
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "productId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&listModel.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		entryController.PostEntry(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Valid request (expect 201)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/v1/shoppinglistentries/1/3",
			strings.NewReader(`{"count": 2, "note": "Test entry", "checked": false}`))
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "productId", "3")
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&listModel.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		mockEntryRepository.EXPECT().Create(&entryModel.UserShoppingListEntry{
			ShoppingListId: 1,
			ProductId:      3,
			Count:          2,
			Note:           "Test entry",
			Checked:        false,
		}).Return(&entryModel.UserShoppingListEntry{
			ShoppingListId: 1,
			ProductId:      3,
			Count:          2,
			Note:           "Test entry",
			Checked:        false,
		}, nil)

		entryController.PostEntry(writer, request)

		assert.Equal(t, http.StatusCreated, writer.Code)
	})

	t.Run("Valid request as admin (expect 201)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/v1/shoppinglistentries/1/3",
			strings.NewReader(`{"count": 2, "note": "Test entry", "checked": false}`))
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "productId", "3")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&listModel.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		mockEntryRepository.EXPECT().Create(&entryModel.UserShoppingListEntry{
			ShoppingListId: 1,
			ProductId:      3,
			Count:          2,
			Note:           "Test entry",
			Checked:        false,
		}).Return(&entryModel.UserShoppingListEntry{
			ShoppingListId: 1,
			ProductId:      3,
			Count:          2,
			Note:           "Test entry",
			Checked:        false,
		}, nil)

		entryController.PostEntry(writer, request)

		assert.Equal(t, http.StatusCreated, writer.Code)
	})

	t.Run("Invalid request, non-numeric listId  (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/v1/shoppinglistentries/abc/3",
			strings.NewReader(`{"count": 2, "note": "Test entry", "checked": false}`))
		ctx := context.WithValue(request.Context(), "listId", "abc")
		ctx = context.WithValue(ctx, "productId", "3")
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		request = request.WithContext(ctx)

		entryController.PostEntry(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})
}

func TestDefaultController_PutEntry(t *testing.T) {
	var mockEntryRepository = userShoppingListEntryMock.NewMockRepository(t)
	var entryRepository Repository = mockEntryRepository

	var mockListRepository = userShoppingListMock.NewMockRepository(t)
	var listRepository userShoppingList.Repository = mockListRepository

	entryController := NewDefaultController(entryRepository, listRepository)

	t.Run("Unauthorized (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/shoppinglistentries/1/2",
			strings.NewReader(`{"count": 10, "note": "Updated entry", "checked": false}`))
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "productId", "2")
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&listModel.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		entryController.PutEntry(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Unauthorized, not matching userIds (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/shoppinglistentries/1/2",
			strings.NewReader(`{"count": 10, "note": "UpdatedUnauthorized entry", "checked": false}`))
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "productId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&listModel.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		entryController.PutEntry(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Valid request (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/shoppinglistentries/1/2",
			strings.NewReader(`{"count": 10, "note": "Updated entry", "checked": false}`))
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "productId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&listModel.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		mockEntryRepository.EXPECT().Update(&entryModel.UserShoppingListEntry{
			ShoppingListId: 1,
			ProductId:      2,
			Count:          10,
			Note:           "Updated entry",
			Checked:        false,
		}).Return(&entryModel.UserShoppingListEntry{
			ShoppingListId: 1,
			ProductId:      2,
			Count:          10,
			Note:           "Updated entry",
			Checked:        false,
		}, nil)

		entryController.PutEntry(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Valid request as admin (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/shoppinglistentries/1/2",
			strings.NewReader(`{"count": 10, "note": "Updated entry", "checked": false}`))
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "productId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&listModel.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		mockEntryRepository.EXPECT().Update(&entryModel.UserShoppingListEntry{
			ShoppingListId: 1,
			ProductId:      2,
			Count:          10,
			Note:           "Updated entry",
			Checked:        false,
		}).Return(&entryModel.UserShoppingListEntry{
			ShoppingListId: 1,
			ProductId:      2,
			Count:          10,
			Note:           "Updated entry",
			Checked:        false,
		}, nil)

		entryController.PutEntry(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Invalid request, non-numeric listId (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/shoppinglistentries/abc/abc",
			strings.NewReader(`{"count": 10, "note": "Updated entry", "checked": false}`))
		ctx := context.WithValue(request.Context(), "listId", "abc")
		ctx = context.WithValue(ctx, "productId", "abc")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		entryController.PutEntry(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Unknown list (expect 404)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/shoppinglistentries/10/2",
			strings.NewReader(`{"count": 10, "note": "Updated entry", "checked": false}`))
		ctx := context.WithValue(request.Context(), "listId", "10")
		ctx = context.WithValue(ctx, "productId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(10)).Return(nil, errors.New(userShoppingList.ErrorListNotFound))

		entryController.PutEntry(writer, request)

		assert.Equal(t, http.StatusNotFound, writer.Code)
	})

	t.Run("Unknown entry (expect 500)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/shoppinglistentries/1/10",
			strings.NewReader(`{"count": 10, "note": "Updated entry", "checked": false}`))
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "productId", "10")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&listModel.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		mockEntryRepository.EXPECT().Update(&entryModel.UserShoppingListEntry{
			ShoppingListId: 1,
			ProductId:      10,
			Count:          10,
			Note:           "Updated entry",
			Checked:        false,
		}).Return(nil, errors.New(ErrorEntryNotFound))

		entryController.PutEntry(writer, request)

		assert.Equal(t, http.StatusInternalServerError, writer.Code)
	})
}

func TestDefaultController_DeleteEntry(t *testing.T) {
	var mockEntryRepository = userShoppingListEntryMock.NewMockRepository(t)
	var entryRepository Repository = mockEntryRepository

	var mockListRepository = userShoppingListMock.NewMockRepository(t)
	var listRepository userShoppingList.Repository = mockListRepository

	entryController := NewDefaultController(entryRepository, listRepository)

	t.Run("Unauthorized (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/shoppinglistentries/1/2", nil)
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "productId", "2")
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&listModel.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		entryController.DeleteEntry(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Unauthorized, not matching userIds (expect 401)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/shoppinglistentries/1/2", nil)
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "productId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&listModel.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		entryController.DeleteEntry(writer, request)

		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})

	t.Run("Valid request (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/shoppinglistentries/1/2", nil)
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "productId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(1))
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&listModel.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		mockEntryRepository.EXPECT().Delete(&entryModel.UserShoppingListEntry{
			ShoppingListId: 1,
			ProductId:      2,
			Count:          0,
			Note:           "",
			Checked:        false,
		}).Return(nil)

		entryController.DeleteEntry(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Valid request as admin (expect 200)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/shoppinglistentries/1/2", nil)
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "productId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&listModel.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		mockEntryRepository.EXPECT().Delete(&entryModel.UserShoppingListEntry{
			ShoppingListId: 1,
			ProductId:      2,
			Count:          0,
			Note:           "",
			Checked:        false,
		}).Return(nil)

		entryController.DeleteEntry(writer, request)

		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Invalid request, non-numeric listId and userId (expect 400)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/shoppinglistentries/abc/abc", nil)
		ctx := context.WithValue(request.Context(), "listId", "abc")
		ctx = context.WithValue(ctx, "productId", "abc")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		entryController.DeleteEntry(writer, request)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
	})

	t.Run("Unknown list (expect 404)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/shoppinglistentries/10/2", nil)
		ctx := context.WithValue(request.Context(), "listId", "10")
		ctx = context.WithValue(ctx, "productId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(10)).Return(nil, errors.New(userShoppingList.ErrorListNotFound))

		entryController.DeleteEntry(writer, request)

		assert.Equal(t, http.StatusNotFound, writer.Code)
	})

	t.Run("Unknown entry (expect 500)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/shoppinglistentries/1/10", nil)
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "productId", "10")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		ctx = context.WithValue(ctx, "auth_userRole", auth.Administrator)
		request = request.WithContext(ctx)

		mockListRepository.EXPECT().FindById(uint64(1)).Return(&listModel.UserShoppingList{
			Id:          1,
			UserId:      1,
			Description: "Test Shopping List",
			Completed:   false,
		}, nil)

		mockEntryRepository.EXPECT().Delete(&entryModel.UserShoppingListEntry{
			ShoppingListId: 1,
			ProductId:      10,
			Count:          0,
			Note:           "",
			Checked:        false,
		}).Return(errors.New(ErrorEntryNotFound))

		entryController.DeleteEntry(writer, request)

		assert.Equal(t, http.StatusInternalServerError, writer.Code)
	})
}
