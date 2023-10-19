package handler

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/auth"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/crypto"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoginHandler(t *testing.T) {
	type fields struct {
		loginHandler *LoginHandler
	}
	type args struct {
		writer  *httptest.ResponseRecorder
		request *http.Request
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "Valid User",
			fields: fields{
				loginHandler: setupLoginHandler(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: httptest.NewRequest(
					"POST",
					"/api/v1/user/login",
					strings.NewReader(`{"email": "ada.lovelace@gmail.com", "password": "12345"}`),
				),
			},
			expectedStatus:   http.StatusCreated,
			expectedResponse: "",
		},
		{
			name: "Malformed JSON",
			fields: fields{
				loginHandler: setupLoginHandler(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: httptest.NewRequest(
					"POST",
					"/api/v1/user/login",
					strings.NewReader(`{"email": "ada.lovelace@gmail.com", "password": "12345"`),
				),
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "",
		},
		{
			name: "Invalid user data, incorrect Type for Email and Password (expected String)",
			fields: fields{
				loginHandler: setupLoginHandler(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: httptest.NewRequest(
					"POST",
					"/api/v1/user/login",
					strings.NewReader(`{"email": 120, "password": 120`),
				),
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.fields.loginHandler.Login(tt.args.writer, tt.args.request)

			// You can then assert the response status and content, and check against your expectations.
			if tt.args.writer.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, but got %d", tt.expectedStatus, tt.args.writer.Code)
			}

			if tt.expectedResponse != "" {
				actualResponse := tt.args.writer.Body.String()
				if actualResponse != tt.expectedResponse {
					t.Errorf("Expected response: %s, but got: %s", tt.expectedResponse, actualResponse)
				}
			}
		})
	}

	loginHandler := setupLoginHandler()

	writer := httptest.NewRecorder()
	request := httptest.NewRequest(
		"POST",
		"/api/v1/user/login",
		strings.NewReader(`{"email": "ada.lovelace@gmail.com", "password": "12345"}`),
	)

	loginHandler.Login(writer, request)

	res := writer.Result()
	var response map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&response)

	assert.NoError(t, err)
	assert.Equal(t, "token", response["access_token"])
	assert.Equal(t, "Bearer", response["token_type"])
	assert.Equal(t, float64(3600), response["expires_in"])
	assert.Equal(t, http.StatusOK, writer.Code)
}

func setupLoginHandler() *LoginHandler {
	var jwtToken, _ = auth.NewJwtTokenGenerator(
		auth.JwtConfig{SignKey: "-----BEGIN EC PRIVATE KEY-----\nProc-Type: 4,ENCRYPTED\nDEK-Info: AES-128-CBC,55D4E42CC38B1EEB2925A48612E432C0\n\n78ZML84hQvNXuIfWKPkzvPsHrEqSGjSgFBbJBWYiDkPHxK3elEKnk8dfXjq0/QVd\nLAf3jQPNLOS0YnZSgGo7GzCCESBvzAKjTt4j8srdJbyybPniw34dfvZr1MTOXNNb\neRsZSh9gCWMt8CA9kROFSnWF/V1MI+0BGOJFXucPOfc=\n-----END EC PRIVATE KEY-----\n"})

	return NewLoginHandler(setupMockRepository(),
		crypto.NewBcryptHasher(), jwtToken)
}

func setupMockRepository() user.Repository {
	repository := user.NewDemoRepository()
	userSlice := setupDemoProductSlice()
	for _, newUser := range userSlice {
		_, _ = repository.Create(newUser)
	}

	return repository
}

func setupDemoProductSlice() []*model.User {
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
		{
			Id:       2,
			Email:    "alan.turin@gmail.com",
			Password: hashedPassword,
			Name:     "Alan Turing",
			Role:     model.Customer,
		},
	}
}
