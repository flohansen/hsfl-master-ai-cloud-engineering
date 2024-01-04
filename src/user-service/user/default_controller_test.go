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

func TestDefaultController_GetUsersByRole(t *testing.T) {
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
			name: "Unauthorized (expect 401)",
			fields: fields{
				userRepository: setupMockRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest("GET", "/api/v1/user/role/1", nil)
					request = request.WithContext(context.WithValue(request.Context(), "userRole", "1"))
					return request
				}(),
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Bad non-numeric request (expect 400)",
			fields: fields{
				userRepository: setupMockRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest("GET", "/api/v1/user/role/abc", nil)
					ctx := context.WithValue(request.Context(), "auth_userId", 1)
					ctx = context.WithValue(ctx, "userRole", "abc")
					request = request.WithContext(ctx)
					return request
				}(),
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Unknown user role (expect 404)",
			fields: fields{
				userRepository: setupMockRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest("GET", "/api/v1/user/role/10", nil)
					ctx := context.WithValue(request.Context(), "auth_userId", 1)
					ctx = context.WithValue(ctx, "userRole", "10")
					request = request.WithContext(ctx)
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
			controller.GetUsersByRole(tt.args.writer, tt.args.request)
			if tt.args.writer.Code != tt.wantStatus {
				t.Errorf("Expected status code %d, got %d", tt.wantStatus, tt.args.writer.Code)
			}
		})
	}

	t.Run("Successfully get existing users by role (expect 200 and users)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/user/role/1", nil)
		ctx := context.WithValue(request.Context(), "auth_userId", 1)
		ctx = context.WithValue(ctx, "userRole", "1")
		request = request.WithContext(ctx)

		controller := defaultController{
			userRepository: setupMockRepository(),
		}

		// when
		controller.GetUsersByRole(writer, request)

		// then
		if writer.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, writer.Code)
		}

		if writer.Header().Get("Content-Type") != "application/json" {
			t.Errorf("Expected content type %s, got %s",
				"application/json", writer.Header().Get("Content-Type"))
		}

		result := writer.Result()
		var response []model.User
		err := json.NewDecoder(result.Body).Decode(&response)
		if err != nil {
			t.Fatal(err.Error())
		}

		if len(response) != 1 {
			t.Errorf("Expected count of user is %d, got %d",
				1, len(response))
		}

		for i, user := range response {
			if user.Role != model.Merchant {
				t.Errorf("Expected role of user %d, got %d", model.Merchant, response[i].Role)
			}
		}
	})
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
					return request.WithContext(context.WithValue(request.Context(), "userId", "abc"))
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
					ctx := context.WithValue(request.Context(), "auth_userId", 10)
					ctx = context.WithValue(ctx, "userId", "10")
					return request.WithContext(ctx)
				}(),
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "Get foreign user as non admin (expect 401)",
			fields: fields{
				userRepository: setupMockRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest("GET", "/api/v1/user/10", nil)
					ctx := context.WithValue(request.Context(), "auth_userId", 1)
					ctx = context.WithValue(ctx, "auth_userRole", 0)
					ctx = context.WithValue(ctx, "userId", "10")
					return request.WithContext(ctx)
				}(),
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Get foreign user as admin (expect 200)",
			fields: fields{
				userRepository: setupMockRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest("GET", "/api/v1/user/1", nil)
					ctx := context.WithValue(request.Context(), "auth_userId", 1)
					ctx = context.WithValue(ctx, "auth_userRole", 2)
					ctx = context.WithValue(ctx, "userId", "2")
					return request.WithContext(ctx)
				}(),
			},
			wantStatus: http.StatusOK,
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

	t.Run("Successfully get user (expect 200 and user)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/user/1", nil)
		ctx := context.WithValue(request.Context(), "auth_userId", 1)
		ctx = context.WithValue(ctx, "auth_userRole", 0)
		ctx = context.WithValue(ctx, "userId", "1")
		request = request.WithContext(ctx)

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

func TestDefaultController_PutUser(t *testing.T) {
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
			name: "Unauthorized (expect 401)",
			fields: fields{
				userRepository: setupMockRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/user/1",
						strings.NewReader(`{"id": 1, "email": "updated@googlemail.com", "name": "Updated user"}`))
					ctx := context.WithValue(request.Context(), "userId", "1")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: "",
		},
		{
			name: "Valid Update as user",
			fields: fields{
				userRepository: setupMockRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/user/1",
						strings.NewReader(`{"id": 1, "email": "updated@googlemail.com", "name": "Updated user"}`))
					ctx := context.WithValue(request.Context(), "auth_userId", 1)
					ctx = context.WithValue(ctx, "userId", "1")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "",
		},
		{
			name: "Valid Update as admin",
			fields: fields{
				userRepository: setupMockRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/user/1",
						strings.NewReader(`{"id": 1, "email": "updated@googlemail.com", "name": "Updated user"}`))
					ctx := context.WithValue(request.Context(), "auth_userId", 1)
					ctx = context.WithValue(ctx, "auth_userRole", 2)
					ctx = context.WithValue(ctx, "userId", "2")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "",
		},
		{
			name: "Valid Update (Partly Fields)",
			fields: fields{
				userRepository: setupMockRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/user/2",
						strings.NewReader(`{"email": "updated@googlemail.com"}`))
					ctx := context.WithValue(request.Context(), "auth_userId", 2)
					ctx = context.WithValue(ctx, "userId", "2")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "",
		},
		{
			name: "Malformed JSON",
			fields: fields{
				userRepository: setupMockRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/user/2",
						strings.NewReader(`{"email": "updated@googlemail.com"`))
					ctx := context.WithValue(request.Context(), "auth_userId", 2)
					ctx = context.WithValue(ctx, "userId", "2")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "",
		},
		{
			name:   "Incorrect Type for email (Non-numeric)",
			fields: fields{},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/user/2",
						strings.NewReader(`{"email": 1234`))
					ctx := context.WithValue(request.Context(), "auth_userId", 2)
					ctx = context.WithValue(ctx, "userId", "2")
					return request.WithContext(ctx)
				}(),
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
			controller.PutUser(tt.args.writer, tt.args.request)

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

func TestDefaultController_DeleteUser(t *testing.T) {
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
			name: "Unauthorized (expect 401)",
			fields: fields{
				userRepository: setupMockRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest("DELETE", "/api/v1/user/1", nil)
					ctx := context.WithValue(request.Context(), "auth_userId", 1)
					ctx = context.WithValue(ctx, "userId", "1")
					return request.WithContext(ctx)
				}(),
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Successfully delete existing user (expect 200)",
			fields: fields{
				userRepository: setupMockRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest("DELETE", "/api/v1/user/1", nil)
					ctx := context.WithValue(request.Context(), "auth_userId", 1)
					ctx = context.WithValue(ctx, "userId", "1")
					return request.WithContext(ctx)
				}(),
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Successfully delete existing user (expect 200)",
			fields: fields{
				userRepository: setupMockRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest("DELETE", "/api/v1/user/1", nil)
					ctx := context.WithValue(request.Context(), "auth_userId", 1)
					ctx = context.WithValue(ctx, "userId", "1")
					return request.WithContext(ctx)
				}(),
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Bad non-numeric request (expect 400)",
			fields: fields{
				userRepository: setupMockRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest("DELETE", "/api/v1/user/abc", nil)
					request = request.WithContext(context.WithValue(request.Context(), "userId", "abc"))
					return request
				}(),
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Unknown user to delete (expect 500)",
			fields: fields{
				userRepository: setupMockRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest("DELETE", "/api/v1/userId/5", nil)
					ctx := context.WithValue(request.Context(), "auth_userId", 1)
					ctx = context.WithValue(ctx, "auth_userRole", 2) // as admin
					ctx = context.WithValue(ctx, "userId", "5")
					return request.WithContext(ctx)
				}(),
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := defaultController{
				userRepository: tt.fields.userRepository,
			}
			controller.DeleteUser(tt.args.writer, tt.args.request)
			if tt.args.writer.Code != tt.wantStatus {
				t.Errorf("Expected status code %d, got %d", tt.wantStatus, tt.args.writer.Code)
			}
		})
	}
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
		{
			Id:       2,
			Email:    "alan.turin@gmail.com",
			Password: hashedPassword,
			Name:     "Alan Turing",
			Role:     model.Merchant,
		},
	}
}

func setupMockRepository() Repository {
	repository := NewDemoRepository()
	usersSlice := setupDemoUserSlice()

	for _, user := range usersSlice {
		repository.Create(user)
	}

	return repository
}
