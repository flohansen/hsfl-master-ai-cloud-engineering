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
	router := New(loginHandler, registerHandler)

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
				assert.Equal(t, http.StatusNotFound, w.Code)
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
				assert.Equal(t, http.StatusNotFound, w.Code)
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
