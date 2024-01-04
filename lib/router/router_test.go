package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	t.Run("should return 404 NOT FOUND if path is unknown", func(t *testing.T) {
		// given
		router := New()
		router.GET("/the/route/without/params", func(w http.ResponseWriter, r *http.Request) {})

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/unknown/route", nil)

		// when
		router.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should route to correct handler without params", func(t *testing.T) {
		// given
		router := New()
		router.GET("/the/route/without/params", func(w http.ResponseWriter, r *http.Request) {})

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/the/route/without/params", nil)

		// when
		router.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should route to correct handler with params", func(t *testing.T) {
		// given
		router := New()
		var ctx context.Context
		router.GET("/the/:route/with/:params", func(w http.ResponseWriter, r *http.Request) {
			ctx = r.Context()
		})

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/the/route/with/params", nil)

		// when
		router.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "route", ctx.Value("route"))
		assert.Equal(t, "params", ctx.Value("params"))
	})

	t.Run("should handle simple handler response", func(t *testing.T) {
		// given
		router := New()
		handlerResponse := "Test response"
		router.GET("/test/handler", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(handlerResponse))
		})

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/test/handler", nil)

		// when
		router.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, handlerResponse, w.Body.String())
	})

	t.Run("should handle handler response with middleware", func(t *testing.T) {
		// given
		router := New()
		handlerResponse := "Test response"
		contextKey := "testKey"
		expectedValue := "testValue"
		var ctx context.Context
		router.GET("/test/context",
			func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(handlerResponse))
				w.WriteHeader(http.StatusOK)
				ctx = r.Context()
			},
			func(w http.ResponseWriter, r *http.Request) *http.Request {
				return r.WithContext(context.WithValue(r.Context(), contextKey, expectedValue))
			})

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/test/context", nil)

		// when
		router.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expectedValue, ctx.Value(contextKey).(string))
		assert.Equal(t, handlerResponse, w.Body.String())
	})
}
