package userShoppingList

import (
	"context"
	"encoding/json"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList/model"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestNewDefaultController(t *testing.T) {
	type args struct {
		userShoppingListRepository Repository
	}
	tests := []struct {
		name string
		args args
		want *defaultController
	}{
		{
			name: "Test construction with DemoRepository",
			args: args{userShoppingListRepository: NewDemoRepository()},
			want: &defaultController{userShoppingListRepository: NewDemoRepository()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDefaultController(tt.args.userShoppingListRepository); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDefaultController() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultController_GetList(t *testing.T) {
	t.Run("Get an existing shopping list (expect 200)", func(t *testing.T) {
		controller := defaultController{
			userShoppingListRepository: setupMockListRepository(),
		}
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglist/1/2", nil)
		ctx := request.Context()
		ctx = context.WithValue(ctx, "listId", "1")
		ctx = context.WithValue(ctx, "userId", "2")
		request = request.WithContext(ctx)

		controller.GetList(writer, request)

		if writer.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, writer.Code)
		}

		var response model.UserShoppingList
		err := json.NewDecoder(writer.Body).Decode(&response)
		if err != nil {
			t.Fatal(err.Error())
		}

		if response.Id != 1 {
			t.Errorf("Expected id of shopping list %d, got %d", 1, response.Id)
		}
		if response.UserId != 2 {
			t.Errorf("Expected user id of shopping list %d, got %d", 2, response.UserId)
		}
	})

	t.Run("Bad non-numeric list ID (expect 400)", func(t *testing.T) {
		controller := defaultController{
			userShoppingListRepository: setupMockListRepository(),
		}
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglist/1/2", nil)
		ctx := request.Context()
		ctx = context.WithValue(ctx, "listId", "abcd")
		ctx = context.WithValue(ctx, "userId", "2")
		request = request.WithContext(ctx)
		controller.GetList(writer, request)

		if writer.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, writer.Code)
		}
	})

	t.Run("Unknown shopping list (expect 404)", func(t *testing.T) {
		controller := defaultController{
			userShoppingListRepository: setupMockListRepository(),
		}
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglist/1/2", nil)
		ctx := request.Context()
		ctx = context.WithValue(ctx, "listId", "5")
		ctx = context.WithValue(ctx, "userId", "2")
		request = request.WithContext(ctx)
		controller.GetList(writer, request)

		if writer.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, writer.Code)
		}
	})
}

func TestDefaultController_PutList(t *testing.T) {
	t.Run("Update an existing shopping list (expect 200)", func(t *testing.T) {
		controller := defaultController{
			userShoppingListRepository: setupMockListRepository(),
		}
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/shoppinglist/1/2", strings.NewReader(`{"checked": true}`))
		request = request.WithContext(context.WithValue(request.Context(), "listId", "1"))
		request = request.WithContext(context.WithValue(request.Context(), "userId", "2"))

		controller.PutList(writer, request)

		if writer.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, writer.Code)
		}
	})

	t.Run("Bad non-numeric list ID (expect 400)", func(t *testing.T) {
		controller := defaultController{
			userShoppingListRepository: setupMockListRepository(),
		}
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/shoppinglist/abc/2", nil)
		request = request.WithContext(context.WithValue(request.Context(), "listId", "abc"))
		request = request.WithContext(context.WithValue(request.Context(), "userId", "2"))

		controller.PutList(writer, request)

		if writer.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, writer.Code)
		}
	})

	t.Run("Malformed JSON request (expect 400)", func(t *testing.T) {
		controller := defaultController{
			userShoppingListRepository: setupMockListRepository(),
		}
		writer := httptest.NewRecorder()

		// Create a malformed JSON request by missing a closing brace '}'.
		request := httptest.NewRequest("PUT", "/api/v1/shoppinglist/1/2", strings.NewReader(`{"checked": true`))
		request = request.WithContext(context.WithValue(request.Context(), "listId", "1"))
		request = request.WithContext(context.WithValue(request.Context(), "userId", "2"))

		controller.PutList(writer, request)

		if writer.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, writer.Code)
		}
	})

	t.Run("Unknown shopping list (expect 500)", func(t *testing.T) {
		controller := defaultController{
			userShoppingListRepository: setupMockListRepositoryError(),
		}
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/api/v1/shoppinglist/5/2", strings.NewReader(`{"checked": true}`))
		request = request.WithContext(context.WithValue(request.Context(), "listId", "5"))
		request = request.WithContext(context.WithValue(request.Context(), "userId", "2"))

		controller.PutList(writer, request)

		if writer.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, writer.Code)
		}
	})
}

func TestDefaultController_PostList(t *testing.T) {
	t.Run("Create a new shopping list (expect 201)", func(t *testing.T) {
		controller := defaultController{
			userShoppingListRepository: setupMockListRepository(),
		}
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/v1/shoppinglist/2", nil)
		request = request.WithContext(context.WithValue(request.Context(), "userId", "2"))

		controller.PostList(writer, request)

		if writer.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, writer.Code)
		}
	})

	t.Run("Bad non-numeric user ID (expect 400)", func(t *testing.T) {
		controller := defaultController{
			userShoppingListRepository: setupMockListRepository(),
		}
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/v1/shoppinglist/2", nil)
		request = request.WithContext(context.WithValue(request.Context(), "userId", "abc"))

		controller.PostList(writer, request)

		if writer.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, writer.Code)
		}
	})
}

func TestDefaultController_GetLists(t *testing.T) {
	t.Run("Retrieve shopping lists for a user (expect 200)", func(t *testing.T) {
		controller := defaultController{
			userShoppingListRepository: setupMockListRepository(),
		}
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglists/2", nil)
		request = request.WithContext(context.WithValue(request.Context(), "userId", "2"))

		controller.GetLists(writer, request)

		if writer.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, writer.Code)
		}

		var response []model.UserShoppingList
		err := json.NewDecoder(writer.Body).Decode(&response)
		if err != nil {
			t.Fatal(err.Error())
		}
	})

	t.Run("Bad non-numeric user ID (expect 400)", func(t *testing.T) {
		controller := defaultController{
			userShoppingListRepository: setupMockListRepository(),
		}
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglists/abc", nil)
		request = request.WithContext(context.WithValue(request.Context(), "userId", "abc"))

		controller.GetLists(writer, request)

		if writer.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, writer.Code)
		}
	})

	t.Run("Internal server error (repository error) (expect 500)", func(t *testing.T) {
		controller := defaultController{
			userShoppingListRepository: setupMockListRepositoryError(),
		}
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/shoppinglists/2", nil)
		request = request.WithContext(context.WithValue(request.Context(), "userId", "2"))

		controller.GetLists(writer, request)

		if writer.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, writer.Code)
		}
	})
}

func TestDefaultController_DeleteList(t *testing.T) {
	t.Run("Delete an existing shopping list (expect 200)", func(t *testing.T) {
		controller := defaultController{
			userShoppingListRepository: setupMockListRepository(),
		}
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/shoppinglist/1/2", nil)
		request = request.WithContext(context.WithValue(request.Context(), "listId", "1"))
		request = request.WithContext(context.WithValue(request.Context(), "userId", "2"))

		controller.DeleteList(writer, request)

		if writer.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, writer.Code)
		}
	})

	t.Run("Bad non-numeric list ID (expect 400)", func(t *testing.T) {
		controller := defaultController{
			userShoppingListRepository: setupMockListRepository(),
		}
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/shoppinglist/abc/2", nil)
		request = request.WithContext(context.WithValue(request.Context(), "listId", "abc"))
		request = request.WithContext(context.WithValue(request.Context(), "userId", "2"))

		controller.DeleteList(writer, request)

		if writer.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, writer.Code)
		}
	})

	t.Run("Unknown shopping list (expect 500)", func(t *testing.T) {
		controller := defaultController{
			userShoppingListRepository: setupMockListRepositoryError(),
		}
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/api/v1/shoppinglist/5/2", nil)
		request = request.WithContext(context.WithValue(request.Context(), "listId", "5"))
		request = request.WithContext(context.WithValue(request.Context(), "userId", "2"))

		controller.DeleteList(writer, request)

		if writer.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, writer.Code)
		}
	})
}

func setupMockListRepository() Repository {
	repository := NewDemoRepository()
	lists := setupDemoListSlice()
	for _, list := range lists {
		repository.Create(list)
	}
	return repository
}

func setupMockListRepositoryError() Repository {
	return &DemoRepository{
		shoppinglists: map[uint64]*model.UserShoppingList{},
	}
}

func setupDemoListSlice() []*model.UserShoppingList {
	return []*model.UserShoppingList{
		{
			Id:          1,
			UserId:      2,
			Description: "Frühstück mit Anne",
			Completed:   false,
		},
		{
			Id:          3,
			UserId:      4,
			Description: "Geburtstagskuchen",
			Completed:   true,
		},
	}
}
