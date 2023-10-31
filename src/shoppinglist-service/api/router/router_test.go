package router

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList"
	shoppingListModel "hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList/model"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingListEntry"
	entryModel "hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingListEntry/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouter(t *testing.T) {
	listRepo := setupMockListRepository()
	entryRepo := setupMockEntryRepository()
	shoppingListController := userShoppingList.NewDefaultController(listRepo)
	shoppingListEntryController := userShoppingListEntry.NewDefaultController(entryRepo)
	router := New(shoppingListController, shoppingListEntryController)

	t.Run("should return 404 NOT FOUND if path is unknown", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/unknown/route", nil)

		// when
		router.ServeHTTP(w, r)

		// then
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("shoppinglist routes", func(t *testing.T) {
		t.Run("GET /api/v1/shoppinglist/:userId should call GetLists", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/shoppinglist/1", nil)

			// when
			router.ServeHTTP(w, r)

			// then
			if w.Code != http.StatusOK {
				t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
			}
		})

		t.Run("GET /api/v1/shoppinglist/:listId/:userId should call GetList", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/shoppinglist/1/2", nil)

			// when
			router.ServeHTTP(w, r)

			// then
			if w.Code != http.StatusOK {
				t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
			}
		})

		t.Run("PUT /api/v1/shoppinglist/:listId/:userId should call PutList", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			jsonRequest := `{"name": "Updated List Name"}`
			r := httptest.NewRequest("PUT", "/api/v1/shoppinglist/1/2", strings.NewReader(jsonRequest))

			// when
			router.ServeHTTP(w, r)

			// then
			if w.Code != http.StatusOK {
				t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
			}
		})

		t.Run("POST /api/v1/shoppinglist/:userId should call PostList", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			jsonRequest := `{"name": "New Shopping List"}`
			r := httptest.NewRequest("POST", "/api/v1/shoppinglist/3", strings.NewReader(jsonRequest))

			// when
			router.ServeHTTP(w, r)

			// then
			if w.Code != http.StatusCreated {
				t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
			}
		})

		t.Run("DELETE /api/v1/shoppinglist/:listId should call DeleteList", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/shoppinglist/1", nil)

			// when
			router.ServeHTTP(w, r)

			// then
			if w.Code != http.StatusOK {
				t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
			}
		})
	})

	t.Run("shoppinglist entries routes", func(t *testing.T) {
		t.Run("GET /api/v1/shoppinglistentries/:listId should call GetEntries", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/shoppinglistentries/1", nil)

			// when
			router.ServeHTTP(w, r)

			// then
			if w.Code != http.StatusOK {
				t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
			}
		})

		t.Run("GET /api/v1/shoppinglistentries/:listId/:productId should call GetEntry", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/shoppinglistentries/1/2", nil)

			// when
			router.ServeHTTP(w, r)

			// then
			if w.Code != http.StatusOK {
				t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
			}
		})

		t.Run("PUT /api/v1/shoppinglistentries/:listId/:productId should call PutEntry", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			jsonRequest := `{"count": 2, "note": "Test entry", "checked": false}`
			r := httptest.NewRequest("PUT", "/api/v1/shoppinglistentries/1/2", strings.NewReader(jsonRequest))

			// when
			router.ServeHTTP(w, r)

			// then
			if w.Code != http.StatusOK {
				t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
			}
		})

		t.Run("POST /api/v1/shoppinglistentries/:listId/:productId should call PostEntry", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			jsonRequest := `{"count": 2, "note": "Test entry", "checked": false}`
			r := httptest.NewRequest("POST", "/api/v1/shoppinglistentries/1/4", strings.NewReader(jsonRequest))

			// when
			router.ServeHTTP(w, r)

			// then
			if w.Code != http.StatusCreated {
				t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
			}
		})

		t.Run("DELETE /api/v1/shoppinglistentries/:listId/:productId should call DeleteEntry", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/shoppinglistentries/1/2", nil)

			// when
			router.ServeHTTP(w, r)

			// then
			if w.Code != http.StatusOK {
				t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
			}
		})
	})

}

func setupMockEntryRepository() userShoppingListEntry.Repository {
	repository := userShoppingListEntry.NewDemoRepository()
	entries := setupDemoEntrySlice()
	for _, entry := range entries {
		repository.Create(entry)
	}
	return repository
}

func setupMockListRepository() userShoppingList.Repository {
	repository := userShoppingList.NewDemoRepository()
	lists := setupDemoListSlice()
	for _, list := range lists {
		repository.Create(list)
	}
	return repository
}

func setupDemoListSlice() []*shoppingListModel.UserShoppingList {
	return []*shoppingListModel.UserShoppingList{
		{
			Id:        1,
			UserId:    2,
			Completed: false,
		},
		{
			Id:        2,
			UserId:    1,
			Completed: true,
		},
		{
			Id:        3,
			UserId:    3,
			Completed: true,
		},
	}
}
func setupDemoEntrySlice() []*entryModel.UserShoppingListEntry {
	return []*entryModel.UserShoppingListEntry{
		{ShoppingListId: 1, ProductId: 1, Count: 3, Note: "Sample entry 1", Checked: false},
		{ShoppingListId: 1, ProductId: 2, Count: 2, Note: "Sample entry 2", Checked: true},
		{ShoppingListId: 2, ProductId: 1, Count: 1, Note: "Sample entry 3", Checked: false},
		{ShoppingListId: 2, ProductId: 2, Count: 4, Note: "Sample entry 3", Checked: false},
	}
}
