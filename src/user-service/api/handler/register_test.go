package handler

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/crypto"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterHandler(t *testing.T) {
	type fields struct {
		registerHandler *RegisterHandler
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
				registerHandler: setUpRegisterHandler(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: httptest.NewRequest(
					"POST",
					"/api/v1/user/register",
					strings.NewReader(`{"email": "ada.lovelace@gmail.com", "password": "123456"}`),
				),
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "",
		},
		{
			name: "Malformed JSON",
			fields: fields{
				registerHandler: setUpRegisterHandler(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: httptest.NewRequest(
					"POST",
					"/api/v1/user/register",
					strings.NewReader(`{"email": "ada.lovelace@gmail.com", "password": "123456"`),
				),
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "",
		},
		{
			name: "Invalid user data, incorrect Type for Email and Password (expected String)",
			fields: fields{
				registerHandler: setUpRegisterHandler(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: httptest.NewRequest(
					"POST",
					"/api/v1/user/register",
					strings.NewReader(`{"email": 120, "password": 120`),
				),
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.fields.registerHandler.Register(tt.args.writer, tt.args.request)

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
		"/api/v1/user/register",
		strings.NewReader(`{"email": "ada.lovelace@gmail.com", "password": "123456"}`),
	)

	loginHandler.Login(writer, request)

	res := writer.Result()
	var response map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&response)

	assert.NoError(t, err)
	assert.Contains(t, response["access_token"], "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9")
	assert.Equal(t, "Bearer", response["token_type"])
	assert.Equal(t, float64(3600), response["expires_in"])
	assert.Equal(t, http.StatusOK, writer.Code)
}

func setUpRegisterHandler() *RegisterHandler {
	return NewRegisterHandler(setupMockRepository(),
		crypto.NewBcryptHasher())
}
