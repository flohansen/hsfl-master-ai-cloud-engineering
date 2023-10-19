package router

import (
	"github.com/stretchr/testify/assert"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/web-service/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	ctrl := service.NewDefaultController()
	router := New(ctrl)

	t.Run("should return shopping list", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/user/shoppingList", nil)
		// when
		router.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should return admin products page", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/admin/products", nil)

		// when
		router.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should merchants products page", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/merchant/products", nil)

		// when
		router.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should return product catalogue page", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/productCatalogue", nil)

		// when
		router.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should return 404 NOT FOUND if target is unknown", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/unknown", nil)

		// when
		router.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
