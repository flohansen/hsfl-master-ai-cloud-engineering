package router

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router/middleware/test"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList"
	userShoppingListMock "hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList/_mock"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingListEntry"
	userShoppingListEntryMock "hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingListEntry/_mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouter(t *testing.T) {
	var mockShoppingListController = userShoppingListMock.NewMockController(t)
	var mockShoppingListEntryController = userShoppingListEntryMock.NewMockController(t)

	authMiddleware := test.CreateEmptyMiddleware()
	var shoppingListController userShoppingList.Controller = mockShoppingListController
	var shoppingListEntryController userShoppingListEntry.Controller = mockShoppingListEntryController

	router := New(&shoppingListController, &shoppingListEntryController, authMiddleware)

	t.Run("/api/v1/shoppinglist/:userId", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/shoppinglist/1", nil)

			mockShoppingListController.EXPECT().GetLists(w, mock.Anything).Run(
				func(_a0 http.ResponseWriter, _a1 *http.Request) {
					_a0.WriteHeader(http.StatusOK)
				})

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call POST handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			jsonRequest := `{"description": "Test List", "checked: false}`
			r := httptest.NewRequest("POST", "/api/v1/shoppinglist/1", strings.NewReader(jsonRequest))

			mockShoppingListController.EXPECT().PostList(w, mock.Anything).Run(
				func(_a0 http.ResponseWriter, _a1 *http.Request) {
					_a0.WriteHeader(http.StatusOK)
				})

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/api/v1/shoppinglist/:userId", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/shoppinglist/1/2", nil)

			mockShoppingListController.EXPECT().GetList(w, mock.Anything).Run(
				func(_a0 http.ResponseWriter, _a1 *http.Request) {
					_a0.WriteHeader(http.StatusOK)
				})

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call PUT handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			jsonRequest := `{"name": "Updated List Name"}`
			r := httptest.NewRequest("PUT", "/api/v1/shoppinglist/1/2", strings.NewReader(jsonRequest))

			mockShoppingListController.EXPECT().PutList(w, mock.Anything).Run(
				func(_a0 http.ResponseWriter, _a1 *http.Request) {
					_a0.WriteHeader(http.StatusOK)
				})

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/api/v1/shoppinglist/:listId", func(t *testing.T) {
		t.Run("should call DELETE handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/shoppinglist/1", nil)

			mockShoppingListController.EXPECT().DeleteList(w, mock.Anything).Run(
				func(_a0 http.ResponseWriter, _a1 *http.Request) {
					_a0.WriteHeader(http.StatusOK)
				})

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/api/v1/shoppinglistentries/:listId", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/shoppinglistentries/1", nil)

			mockShoppingListEntryController.EXPECT().GetEntries(w, mock.Anything).Run(
				func(_a0 http.ResponseWriter, _a1 *http.Request) {
					_a0.WriteHeader(http.StatusOK)
				})

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/api/v1/shoppinglistentries/:listId/:productId", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/shoppinglistentries/1/2", nil)

			mockShoppingListEntryController.EXPECT().GetEntry(w, mock.Anything).Run(
				func(_a0 http.ResponseWriter, _a1 *http.Request) {
					_a0.WriteHeader(http.StatusOK)
				})

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call PUT handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			jsonRequest := `{"count": 10, "note": "Test entry", "checked": false}`
			r := httptest.NewRequest("PUT", "/api/v1/shoppinglistentries/1/2", strings.NewReader(jsonRequest))

			mockShoppingListEntryController.EXPECT().PutEntry(w, mock.Anything).Run(
				func(_a0 http.ResponseWriter, _a1 *http.Request) {
					_a0.WriteHeader(http.StatusOK)
				})

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call POST handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			jsonRequest := `{"count": 2, "note": "Test entry", "checked": false}`
			r := httptest.NewRequest("POST", "/api/v1/shoppinglistentries/1/10", strings.NewReader(jsonRequest))

			mockShoppingListEntryController.EXPECT().PostEntry(w, mock.Anything).Run(
				func(_a0 http.ResponseWriter, _a1 *http.Request) {
					_a0.WriteHeader(http.StatusOK)
				})

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call DELETE handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/shoppinglistentries/1/2", nil)

			mockShoppingListEntryController.EXPECT().DeleteEntry(w, mock.Anything).Run(
				func(_a0 http.ResponseWriter, _a1 *http.Request) {
					_a0.WriteHeader(http.StatusOK)
				})

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
}
