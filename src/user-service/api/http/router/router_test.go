package router

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router/middleware/test"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/api/http/handler"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/auth"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/crypto"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user"
	userMock "hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user/_mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouter(t *testing.T) {
	var mockUserController = userMock.NewMockController(t)
	var mockUserRepository = userMock.NewMockRepository(t)

	var userController user.Controller = mockUserController
	var userRepository user.Repository = mockUserRepository

	loginHandler := setupLoginHandler(userRepository)
	registerHandler := setUpRegisterHandler(userRepository)
	authMiddleware := test.CreateEmptyMiddleware()

	router := New(loginHandler, registerHandler, &userController, authMiddleware)

	t.Run("/api/v1/authentication/login/", func(t *testing.T) {
		t.Run("should call POST handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/authentication/login/", nil)
			r.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	})

	t.Run("/api/v1/authentication/register/", func(t *testing.T) {
		t.Run("should call POST handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/authentication/register/", nil)
			r.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	})

	t.Run("/api/v1/user/role/:userRole", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/user/role/1", nil)

			mockUserController.EXPECT().GetUsersByRole(w, mock.Anything).Run(
				func(_a0 http.ResponseWriter, _a1 *http.Request) {
					_a0.WriteHeader(http.StatusOK)
				})

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/api/v1/user/:userId", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/user/1", nil)

			mockUserController.EXPECT().GetUser(w, mock.Anything).Run(
				func(_a0 http.ResponseWriter, _a1 *http.Request) {
					_a0.WriteHeader(http.StatusOK)
				})

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call PUT handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			jsonRequest := `{"id": 1, "email": "updated@googlemail.com", "name": "Updated user"}`
			r := httptest.NewRequest("PUT", "/api/v1/user/1", strings.NewReader(jsonRequest))

			mockUserController.EXPECT().PutUser(w, mock.Anything).Run(
				func(_a0 http.ResponseWriter, _a1 *http.Request) {
					_a0.WriteHeader(http.StatusOK)
				})

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call DELETE handler", func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/user/1", nil)

			mockUserController.EXPECT().DeleteUser(w, mock.Anything).Run(
				func(_a0 http.ResponseWriter, _a1 *http.Request) {
					_a0.WriteHeader(http.StatusOK)
				})

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
}

func setupLoginHandler(userRepository user.Repository) *handler.LoginHandler {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic("Error generating private key for testing.")
	}

	derFormatKey, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		panic("Error converting ECDSA key to DER format")
	}

	pemKey := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: derFormatKey,
	})

	pemPrivateKey := string(pemKey)

	var jwtToken, _ = auth.NewJwtTokenGenerator(
		auth.JwtConfig{PrivateKey: pemPrivateKey})

	return handler.NewLoginHandler(userRepository,
		crypto.NewBcryptHasher(), jwtToken)
}

func setUpRegisterHandler(userRepository user.Repository) *handler.RegisterHandler {
	return handler.NewRegisterHandler(userRepository,
		crypto.NewBcryptHasher())
}
