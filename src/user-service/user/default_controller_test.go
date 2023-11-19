package user

import (
	"context"
	"encoding/json"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/crypto"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user/model"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestNewDefaultController(t *testing.T) {
	type args struct {
		userRepository Repository
	}
	tests := []struct {
		name string
		args args
		want *defaultController
	}{
		{
			name: "Test construction with DemoRepository",
			args: args{userRepository: NewDemoRepository()},
			want: &defaultController{userRepository: NewDemoRepository()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDefaultController(tt.args.userRepository); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDefaultController() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultController_GetUser(t *testing.T) {
	type fields struct {
		userRepository Repository
	}
	type args struct {
		writer  *httptest.ResponseRecorder
		request *http.Request
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
	}{
		{
			name: "Bad non-numeric request (expect 400)",
			fields: fields{
				userRepository: setupMockRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest("GET", "/api/v1/user/abc", nil)
					request = request.WithContext(context.WithValue(request.Context(), "userId", "abc"))
					return request
				}(),
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Unknown user (expect 404)",
			fields: fields{
				userRepository: setupMockRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest("GET", "/api/v1/user/10", nil)
					request = request.WithContext(context.WithValue(request.Context(), "userId", "10"))
					return request
				}(),
			},
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := defaultController{
				userRepository: tt.fields.userRepository,
			}
			controller.GetUser(tt.args.writer, tt.args.request)
			if tt.args.writer.Code != tt.wantStatus {
				t.Errorf("Expected status code %d, got %d", tt.wantStatus, tt.args.writer.Code)
			}
		})
	}

	t.Run("Successfully get existing user (expect 200 and user)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/user/1", nil)
		request = request.WithContext(context.WithValue(request.Context(), "userId", "1"))

		controller := defaultController{
			userRepository: setupMockRepository(),
		}

		// when
		controller.GetUser(writer, request)

		// then
		if writer.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, writer.Code)
		}

		if writer.Header().Get("Content-Type") != "application/json" {
			t.Errorf("Expected content type %s, got %s",
				"application/json", writer.Header().Get("Content-Type"))
		}

		result := writer.Result()
		var response model.User
		err := json.NewDecoder(result.Body).Decode(&response)
		if err != nil {
			t.Fatal(err.Error())
		}

		if response.Id != 1 {
			t.Errorf("Expected id of user %d, got %d", 1, response.Id)
		}

		if response.Email != "ada.lovelace@gmail.com" {
			t.Errorf("Expected email of user %s, got %s", "ada.lovelace@gmail.com", response.Email)
		}

		if response.Name != "Ada Lovelace" {
			t.Errorf("Expected name of user %s, got %s", "Ada Lovelace", response.Name)
		}

		if response.Role != model.Customer {
			t.Errorf("Got false user role")
		}
	})
}

func TestDefaultController_PostUser(t *testing.T) {
	type fields struct {
		userRepository Repository
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
				userRepository: setupMockRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: httptest.NewRequest(
					"POST",
					"/api/v1/user",
					strings.NewReader(`{"id": 3, "email": "example@googlemail.com", "password": "password", "name": "Example name"}`),
				),
			},
			expectedStatus:   http.StatusCreated,
			expectedResponse: "",
		},
		{
			name: "Valid User (Partly Fields)",
			fields: fields{
				userRepository: setupMockRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: httptest.NewRequest(
					"POST",
					"/api/v1/user",
					strings.NewReader(`{"email": "example@googlemail.com", "password": "password"}`),
				),
			},
			expectedStatus:   http.StatusCreated,
			expectedResponse: "",
		},
		{
			name: "Malformed JSON",
			fields: fields{
				userRepository: setupMockRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: httptest.NewRequest(
					"POST",
					"/api/v1/user",
					strings.NewReader(`{"email": "example@googlemail.com"`),
				),
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "",
		},
		{
			name: "Invalid user, incorrect Type for email (Non-numeric)",
			fields: fields{
				userRepository: setupMockRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: httptest.NewRequest(
					"POST",
					"/api/v1/user",
					strings.NewReader(`{"id": 3, "email": 1234, "password": "password", "name": "Example name"}`),
				),
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := defaultController{
				userRepository: tt.fields.userRepository,
			}
			controller.PostUser(tt.args.writer, tt.args.request)

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

func setupMockRepository() Repository {
	bcryptHasher := crypto.NewBcryptHasher()
	hashedPassword, _ := bcryptHasher.Hash([]byte("123456"))

	repository := NewDemoRepository()
	usersSlice := []*model.User{
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
	for _, user := range usersSlice {
		repository.Create(user)
	}

	return repository
}
