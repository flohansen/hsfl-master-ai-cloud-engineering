package router

import (
	"github.com/stretchr/testify/assert"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/api/handler"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/auth"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/crypto"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouter(t *testing.T) {
	loginHandler := setupLoginHandler()
	registerHandler := setUpRegisterHandler()

	userRepo := setupUserRepository()
	userController := user.NewDefaultController(userRepo)
	router := New(loginHandler, registerHandler, userController)

	t.Run("should return 404 NOT FOUND if path is unknown", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/unknown/route", nil)

		// when
		router.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("/api/v1/user/login", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not POST", func(t *testing.T) {
			tests := []string{"HEAD", "CONNECT", "OPTIONS", "TRACE", "PATCH", "GET", "DELETE", "PUT"}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/api/v1/user/login", nil)

				// when
				router.ServeHTTP(w, r)

				// then
				if test == "GET" || test == "PUT" || test == "DELETE" {
					assert.Equal(t, http.StatusBadRequest, w.Code)
				} else {
					assert.Equal(t, http.StatusNotFound, w.Code)
				}
			}
		})

		t.Run("should call POST handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			jsonRequest := `{"email": "ada.lovelace@gmail.com", "password": "123456"}`
			r := httptest.NewRequest("POST", "/api/v1/user/login", strings.NewReader(jsonRequest))

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/api/v1/user/register", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not POST", func(t *testing.T) {
			tests := []string{"HEAD", "CONNECT", "OPTIONS", "TRACE", "PATCH", "GET", "DELETE", "PUT"}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/api/v1/user/register", nil)

				// when
				router.ServeHTTP(w, r)

				// then
				if test == "GET" || test == "PUT" || test == "DELETE" {
					assert.Equal(t, http.StatusBadRequest, w.Code)
				} else {
					assert.Equal(t, http.StatusNotFound, w.Code)
				}
			}
		})

		t.Run("should call POST handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/user/register", strings.NewReader(`{"email": "grace.hopper@gmail.com", "password": "123456", "name": "Grace Hopper", "role": 0}`))

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/api/v1/user", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not GET", func(t *testing.T) {
			tests := []string{"DELETE", "PUT", "HEAD", "CONNECT", "OPTIONS", "TRACE", "PATCH"}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/api/v1/products", nil)

				// when
				router.ServeHTTP(w, r)

				// then
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		})

		t.Run("should call GET handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/user/", nil)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call POST handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			jsonRequest := `{"id": 3, "email": "example@googlemail.com", "password": "password", "name": "Example name"}`
			r := httptest.NewRequest("POST", "/api/v1/user/", strings.NewReader(jsonRequest))

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusCreated, w.Code)
		})
	})

	t.Run("/api/v1/user/:userId", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not GET, DELETE or PUT", func(t *testing.T) {
			tests := []string{"HEAD", "CONNECT", "OPTIONS", "TRACE", "PATCH", "POST"}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/api/v1/user/1", nil)

				// when
				router.ServeHTTP(w, r)

				// then
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		})

		t.Run("should call GET handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/user/1", nil)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call PUT handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			jsonRequest := `{"id": 1, "email": "updated@googlemail.com", "name": "Updated user"}`
			r := httptest.NewRequest("PUT", "/api/v1/user/1", strings.NewReader(jsonRequest))

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call DELETE handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/user/1", nil)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
}

func setupLoginHandler() *handler.LoginHandler {
	var jwtToken, _ = auth.NewJwtTokenGenerator(
		auth.JwtConfig{SignKey: "../../auth/test-token"})

	return handler.NewLoginHandler(setupMockRepository(),
		crypto.NewBcryptHasher(), jwtToken)
}

func setUpRegisterHandler() *handler.RegisterHandler {
	return handler.NewRegisterHandler(setupMockRepository(),
		crypto.NewBcryptHasher())
}

func setupMockRepository() user.Repository {
	repository := user.NewDemoRepository()
	userSlice := setupDemoUserSlice()
	for _, newUser := range userSlice {
		_, _ = repository.Create(newUser)
	}

	return repository
}

func setupUserRepository() user.Repository {
	repository := user.NewDemoRepository()
	userSlice := setupDemoUserSlice()
	for _, u := range userSlice {
		_, err := repository.Create(u)
		if err != nil {
			return nil
		}
	}

	return repository
}

func setupDemoUserSlice() []*model.User {
	bcryptHasher := crypto.NewBcryptHasher()
	hashedPassword, _ := bcryptHasher.Hash([]byte("123456"))

	return []*model.User{
		{
			Id:       1,
			Email:    "ada.lovelace@gmail.com",
			Password: hashedPassword,
			Name:     "Ada Lovelace",
			Role:     model.Customer,
		},
	}
}
