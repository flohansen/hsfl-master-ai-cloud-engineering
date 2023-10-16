package router

import (
	"context"
	mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	ctrl := gomock.NewController(t)

	userController := mocks.NewMockController(ctrl)
	router := New(userController)

	t.Run("/api/v1/users", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not GET", func(t *testing.T) {
			tests := []string{"HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}
			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/api/v1/users", nil)

				// when
				router.ServeHTTP(w, r)

				// then
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		})

		t.Run("should call GET handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/users", nil)

			userController.
				EXPECT().
				GetUsers(w, r).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/api/v1/login", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not POST", func(t *testing.T) {

			tests := []string{"GET", "HEAD", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/api/v1/login", nil)

				// when
				router.ServeHTTP(w, r)

				// then
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		})

		t.Run("should call POST handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/login", nil)

			userController.
				EXPECT().
				Login(w, r).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/api/v1/register", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not POST", func(t *testing.T) {
			tests := []string{"GET", "HEAD", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/api/v1/register", nil)

				// when
				router.ServeHTTP(w, r)

				// then
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		})

		t.Run("should call POST handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/register", nil)

			userController.
				EXPECT().
				Register(w, r).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/api/v1/users/:username", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not GET, DELETE or PUT", func(t *testing.T) {
			tests := []string{"HEAD", "POST", "CONNECT", "OPTIONS", "TRACE", "PATCH"}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/api/v1/users/tester", nil)

				// when
				router.ServeHTTP(w, r)

				// then
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		})

		t.Run("should call GET handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/users/tester", nil)

			userController.
				EXPECT().
				GetUser(w, r.WithContext(context.WithValue(r.Context(), "username", "tester"))).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call PUT handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/users/tester", nil)

			userController.
				EXPECT().
				PutUser(w, r.WithContext(context.WithValue(r.Context(), "username", "tester"))).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call DELETE handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/users/tester", nil)

			userController.
				EXPECT().
				DeleteUser(w, r.WithContext(context.WithValue(r.Context(), "username", "tester"))).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
}
