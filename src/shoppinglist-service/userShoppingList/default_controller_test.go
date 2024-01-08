package userShoppingList

import (
	"context"
	"encoding/json"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList/model"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestNewDefaultController(t *testing.T) {
	type args struct {
		userShoppingListRepository Repository
	}
	tests := []struct {
		name string
		args args
		want *DefaultController
	}{
		{
			name: "Test construction with DemoRepository",
			args: args{userShoppingListRepository: NewDemoRepository()},
			want: &DefaultController{userShoppingListRepository: NewDemoRepository()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDefaultController(tt.args.userShoppingListRepository); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDefaultController() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultController_GetLists(t *testing.T) {
	type fields struct {
		userShoppingListRepository Repository
	}
	type args struct {
		writer  *httptest.ResponseRecorder
		request *http.Request
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		expectedStatus int
	}{
		{
			name: "Unauthorized (expect 401)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"GET",
						"/api/v1/shoppinglist/2",
						nil)
					ctx := context.WithValue(request.Context(), "userId", "2")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Unauthorized, not matching userIds (expect 401)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"GET",
						"/api/v1/shoppinglist/2",
						nil)
					ctx := context.WithValue(request.Context(), "userId", "2")
					ctx = context.WithValue(ctx, "auth_userId", uint64(1))
					ctx = context.WithValue(ctx, "auth_userRole", int64(1))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Valid request (expect 200)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"GET",
						"/api/v1/shoppinglist/2",
						nil)
					ctx := context.WithValue(request.Context(), "userId", "2")
					ctx = context.WithValue(ctx, "auth_userId", uint64(2))
					ctx = context.WithValue(ctx, "auth_userRole", int64(1))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Valid request as admin (expect 200)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"GET",
						"/api/v1/shoppinglist/2",
						nil)
					ctx := context.WithValue(request.Context(), "userId", "2")
					ctx = context.WithValue(ctx, "auth_userId", uint64(3))
					ctx = context.WithValue(ctx, "auth_userRole", int64(2))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid request, non-numeric userId (expect 400)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"GET",
						"/api/v1/shoppinglist/abc",
						nil)
					ctx := context.WithValue(request.Context(), "userId", "abc")
					ctx = context.WithValue(ctx, "auth_userId", uint64(2))
					ctx = context.WithValue(ctx, "auth_userRole", int64(1))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := DefaultController{
				userShoppingListRepository: setupMockListRepository(),
			}
			controller.GetLists(tt.args.writer, tt.args.request)
			if tt.args.writer.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, tt.args.writer.Code)
			}
		})
	}

	t.Run("Successfully get user shopping lists (expect 200 and shopping lists)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		var request = httptest.NewRequest(
			"GET",
			"/api/v1/shoppinglist/2",
			nil)
		ctx := context.WithValue(request.Context(), "userId", "2")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		ctx = context.WithValue(ctx, "auth_userRole", int64(1))
		request = request.WithContext(ctx)

		controller := DefaultController{
			userShoppingListRepository: setupMockListRepository(),
		}

		controller.GetLists(writer, request)

		if writer.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, writer.Code)
		}

		if writer.Header().Get("Content-Type") != "application/json" {
			t.Errorf("Expected content type %s, got %s",
				"application/json", writer.Header().Get("Content-Type"))
		}

		result := writer.Result()
		var response []model.UserShoppingList
		err := json.NewDecoder(result.Body).Decode(&response)
		if err != nil {
			t.Fatal(err.Error())
		}

		if len(response) != 2 {
			t.Errorf("Expected count of shopping lists is %d, got %d",
				2, len(response))
		}

		for i, list := range response {
			if list.UserId != response[i].UserId {
				t.Errorf("Expected userId of shopping lists is %d, got %d", list.UserId, response[i].UserId)
			}
		}
	})
}

func TestDefaultController_GetList(t *testing.T) {
	type fields struct {
		userShoppingListRepository Repository
	}
	type args struct {
		writer  *httptest.ResponseRecorder
		request *http.Request
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		expectedStatus int
	}{
		{
			name: "Unauthorized (expect 401)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"GET",
						"/api/v1/shoppinglist/1/2",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "userId", "2")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Unauthorized, not matching userIds (expect 401)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"GET",
						"/api/v1/shoppinglist/1/2",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "userId", "2")
					ctx = context.WithValue(ctx, "auth_userId", uint64(1))
					ctx = context.WithValue(ctx, "auth_userRole", int64(1))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Valid request (expect 200)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"GET",
						"/api/v1/shoppinglist/1/2",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "userId", "2")
					ctx = context.WithValue(ctx, "auth_userId", uint64(2))
					ctx = context.WithValue(ctx, "auth_userRole", int64(1))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Valid request as admin (expect 200)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"GET",
						"/api/v1/shoppinglist/1/2",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "userId", "2")
					ctx = context.WithValue(ctx, "auth_userId", uint64(3))
					ctx = context.WithValue(ctx, "auth_userRole", int64(2))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid request, non-numeric userId (expect 400)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"GET",
						"/api/v1/shoppinglist/abc/abc",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "ab")
					ctx = context.WithValue(ctx, "userId", "abc")
					ctx = context.WithValue(ctx, "auth_userId", "abc")
					ctx = context.WithValue(ctx, "auth_userRole", "abc")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Unknown shopping list (expect 404)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"GET",
						"/api/v1/shoppinglist/4/2",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "4")
					ctx = context.WithValue(ctx, "userId", "2")
					ctx = context.WithValue(ctx, "auth_userId", uint64(2))
					ctx = context.WithValue(ctx, "auth_userRole", int64(0))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := DefaultController{
				userShoppingListRepository: setupMockListRepository(),
			}
			controller.GetList(tt.args.writer, tt.args.request)
			if tt.args.writer.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, tt.args.writer.Code)
			}
		})
	}
}

func TestDefaultController_PutList(t *testing.T) {
	type fields struct {
		userShoppingListRepository Repository
	}
	type args struct {
		writer  *httptest.ResponseRecorder
		request *http.Request
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		expectedStatus int
	}{
		{
			name: "Unauthorized (expect 401)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/shoppinglist/1/2",
						strings.NewReader(`{"description": "Updated list"}`))
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "userId", "2")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Unauthorized, not matching userIds (expect 401)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/shoppinglist/1/2",
						strings.NewReader(`{"description": "Updated list"}`))
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "userId", "2")
					ctx = context.WithValue(ctx, "auth_userId", uint64(1))
					ctx = context.WithValue(ctx, "auth_userRole", int64(1))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Valid request (expect 200)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/shoppinglist/1/2",
						strings.NewReader(`{"description": "Updated list"}`))
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "userId", "2")
					ctx = context.WithValue(ctx, "auth_userId", uint64(2))
					ctx = context.WithValue(ctx, "auth_userRole", int64(1))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Valid request as admin(expect 200)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/shoppinglist/1/2",
						strings.NewReader(`{"description": "Updated list"}`))
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "userId", "2")
					ctx = context.WithValue(ctx, "auth_userId", uint64(1))
					ctx = context.WithValue(ctx, "auth_userRole", int64(2))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid request, non-numeric userId (expect 400)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/shoppinglist/abc/abc",
						strings.NewReader(`{"description": "Updated list"}`))
					ctx := context.WithValue(request.Context(), "listId", "ab")
					ctx = context.WithValue(ctx, "userId", "abc")
					ctx = context.WithValue(ctx, "auth_userId", "abc")
					ctx = context.WithValue(ctx, "auth_userRole", "abc")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid request, malformed JSON (expect 400)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/shoppinglist/1/2",
						strings.NewReader(`{"description": "Updated list"`))
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "userId", "2")
					ctx = context.WithValue(ctx, "auth_userId", uint64(2))
					ctx = context.WithValue(ctx, "auth_userRole", int64(1))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := DefaultController{
				userShoppingListRepository: setupMockListRepository(),
			}
			controller.PutList(tt.args.writer, tt.args.request)
			if tt.args.writer.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, tt.args.writer.Code)
			}
		})
	}
}

func TestDefaultController_PostList(t *testing.T) {
	type fields struct {
		userShoppingListRepository Repository
	}
	type args struct {
		writer  *httptest.ResponseRecorder
		request *http.Request
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		expectedStatus int
	}{
		{
			name: "Unauthorized (expect 401)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"POST",
						"/api/v1/shoppinglist/2",
						strings.NewReader(`{"description": "New list", "checked": false}`))
					ctx := context.WithValue(request.Context(), "userId", "2")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Unauthorized, not matching userIds (expect 401)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"POST",
						"/api/v1/shoppinglist/2",
						strings.NewReader(`{"description": "New list", "checked": false}`))
					ctx := context.WithValue(request.Context(), "userId", "2")
					ctx = context.WithValue(ctx, "auth_userId", uint64(1))
					ctx = context.WithValue(ctx, "auth_userRole", int64(1))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Valid request (expect 201)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"POST",
						"/api/v1/shoppinglist/2",
						strings.NewReader(`{"description": "New list", "checked": false}`))
					ctx := context.WithValue(request.Context(), "userId", "2")
					ctx = context.WithValue(ctx, "auth_userId", uint64(2))
					ctx = context.WithValue(ctx, "auth_userRole", int64(1))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Valid request as admin (expect 201)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"POST",
						"/api/v1/shoppinglist/1",
						strings.NewReader(`{"description": "New list", "checked": false}`))
					ctx := context.WithValue(request.Context(), "userId", "1")
					ctx = context.WithValue(ctx, "auth_userId", uint64(2))
					ctx = context.WithValue(ctx, "auth_userRole", int64(2))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid request, non-numeric userId (expect 400)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"POST",
						"/api/v1/shoppinglist/abc",
						strings.NewReader(`{"description": "New list", "checked": false}`))
					ctx := context.WithValue(request.Context(), "userId", "abc")
					ctx = context.WithValue(ctx, "auth_userId", uint64(2))
					ctx = context.WithValue(ctx, "auth_userRole", int64(2))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid request, malformed JSON (expect 400)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"POST",
						"/api/v1/shoppinglist/2",
						strings.NewReader(`{"description": "New list", "checked": false`))
					ctx := context.WithValue(request.Context(), "userId", "2")
					ctx = context.WithValue(ctx, "auth_userId", uint64(2))
					ctx = context.WithValue(ctx, "auth_userRole", int64(2))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := DefaultController{
				userShoppingListRepository: setupMockListRepository(),
			}
			controller.PostList(tt.args.writer, tt.args.request)
			if tt.args.writer.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, tt.args.writer.Code)
			}
		})
	}
}

func TestDefaultController_Delete(t *testing.T) {
	type fields struct {
		userShoppingListRepository Repository
	}
	type args struct {
		writer  *httptest.ResponseRecorder
		request *http.Request
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		expectedStatus int
	}{
		{
			name: "Unauthorized (expect 401)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"DELETE",
						"/api/v1/shoppinglist/1",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "1")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Unauthorized, not matching userIds (expect 401)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"DELETE",
						"/api/v1/shoppinglist/1",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "auth_userId", uint64(1))
					ctx = context.WithValue(ctx, "auth_userRole", int64(0))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Valid request (expect 200)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"DELETE",
						"/api/v1/shoppinglist/1",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "auth_userId", uint64(2))
					ctx = context.WithValue(ctx, "auth_userRole", int64(0))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Valid request as admin (expect 200)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"DELETE",
						"/api/v1/shoppinglist/2",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "2")
					ctx = context.WithValue(ctx, "auth_userId", uint64(3))
					ctx = context.WithValue(ctx, "auth_userRole", int64(2))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid request, non-numeric listId (expect 400)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"DELETE",
						"/api/v1/shoppinglist/abc",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "abc")
					ctx = context.WithValue(ctx, "auth_userId", uint64(3))
					ctx = context.WithValue(ctx, "auth_userRole", int64(2))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Unknown list (expect 404)",
			fields: fields{
				userShoppingListRepository: setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"DELETE",
						"/api/v1/shoppinglist/10",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "10")
					ctx = context.WithValue(ctx, "auth_userId", uint64(1))
					ctx = context.WithValue(ctx, "auth_userRole", int64(2))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := DefaultController{
				userShoppingListRepository: setupMockListRepository(),
			}
			controller.DeleteList(tt.args.writer, tt.args.request)
			if tt.args.writer.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, tt.args.writer.Code)
			}
		})
	}
}

func setupMockListRepository() Repository {
	repository := NewDemoRepository()
	lists := setupDemoListSlice()
	for _, list := range lists {
		repository.Create(list)
	}
	return repository
}

func setupDemoListSlice() []*model.UserShoppingList {
	return []*model.UserShoppingList{
		{
			Id:          1,
			UserId:      2,
			Description: "Frühstück mit Anne",
			Completed:   false,
		},
		{
			Id:          2,
			UserId:      2,
			Description: "Suppe",
			Completed:   false,
		},
		{
			Id:          3,
			UserId:      4,
			Description: "Geburtstagskuchen",
			Completed:   true,
		},
	}
}
