package router

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
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

	t.Run("should not call middleware", func(t *testing.T) {
		router := New()

		router.USE("/unknown/route", func(w http.ResponseWriter, r *http.Request, next Next) {
			w.WriteHeader(http.StatusAlreadyReported)
			next(r)
		})

		router.GET("/", func(w http.ResponseWriter, r *http.Request) {})

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		// when
		router.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should call middleware without params", func(t *testing.T) {
		router := New()

		router.USE("/the/route/without/params", func(w http.ResponseWriter, r *http.Request, next Next) {
			next(createRequestContext(r, []string{"hello"}, []string{"world"}))
		})

		var ctx context.Context
		router.GET("/the/route/without/params", func(w http.ResponseWriter, r *http.Request) {
			ctx = r.Context()
		})

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/the/route/without/params", nil)

		// when
		router.ServeHTTP(w, r)

		// then
		assert.Equal(t, "world", ctx.Value("hello"))
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should call middleware without params and all sub-dirs", func(t *testing.T) {
		router := New()

		router.USE("/the/route/without/params", func(w http.ResponseWriter, r *http.Request, next Next) {
			w.WriteHeader(http.StatusAlreadyReported)
			next(r)
		})

		router.GET("/the/route/without/params/but/sub/folders/are/allowed", func(w http.ResponseWriter, r *http.Request) {})

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/the/route/without/params/but/sub/folders/are/allowed", nil)

		// when
		router.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusAlreadyReported, w.Code)
	})

	t.Run("should call middleware with params", func(t *testing.T) {
		// given
		router := New()
		var ctx context.Context
		router.USE("/the/:route/with/:params", func(w http.ResponseWriter, r *http.Request, next Next) {
			ctx = r.Context()
			next(r)
		})

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/the/route/with/params", nil)

		// when
		router.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "route", ctx.Value("route"))
		assert.Equal(t, "params", ctx.Value("params"))
	})

	t.Run("don't route if next in middleware is not called", func(t *testing.T) {
		// given
		router := New()
		router.USE("/", func(w http.ResponseWriter, r *http.Request, next Next) {
			w.WriteHeader(http.StatusUnauthorized)
		})

		router.GET("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte("doesn't matter"))
		})

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		// when
		router.ServeHTTP(w, r)
		s, err := io.ReadAll(r.Body)

		// then
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.NoError(t, err)
		assert.Equal(t, s, []byte{})
	})
}
