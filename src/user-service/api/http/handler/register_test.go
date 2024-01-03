package handler

import (
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
					strings.NewReader(`{"email": "grace.hopper@gmail.com", "password": "123456", "name": "Grace Hopper", "role": 0}`),
				),
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "",
		},
		{
			name: "Invalid Request - Empty Password",
			fields: fields{
				registerHandler: setUpRegisterHandler(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: httptest.NewRequest(
					"POST",
					"/api/v1/user/register",
					strings.NewReader(`{"email": "grace.hopper2@gmail.com", "password": "", "name": "Grace Hopper", "role": 0}`),
				),
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "",
		},
		{
			name: "User already exists",
			fields: fields{
				registerHandler: setUpRegisterHandler(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: httptest.NewRequest(
					"POST",
					"/api/v1/user/register",
					strings.NewReader(`{"email": "ada.lovelace@gmail.com", "password": "123456", "name": "Ada Lovelace", "role": 0}`),
				),
			},
			expectedStatus:   http.StatusConflict,
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
					strings.NewReader(`{"email": "grace.hopper@gmail.com", "password": "123456", "name": "Grace Hopper", "role": 0`),
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
					strings.NewReader(`{"email": "grace.hopper@gmail.com", "password": 123456, "name": "Grace Hopper", "role": false}`),
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
}

func setUpRegisterHandler() *RegisterHandler {
	return NewRegisterHandler(setupMockRepository(),
		crypto.NewBcryptHasher())
}
