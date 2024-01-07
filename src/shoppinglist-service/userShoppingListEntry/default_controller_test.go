package userShoppingListEntry

import (
	"context"
	"encoding/json"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList"
	listModel "hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList/model"
	entryModel "hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingListEntry/model"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestNewDefaultController(t *testing.T) {
	type args struct {
		userShoppingListEntryRepository Repository
		userShoppingListRepository      userShoppingList.Repository
	}
	tests := []struct {
		name string
		args args
		want *DefaultController
	}{
		{
			name: "Test construction with DemoRepository",
			args: args{
				userShoppingListEntryRepository: NewDemoRepository(),
				userShoppingListRepository:      userShoppingList.NewDemoRepository(),
			},
			want: &DefaultController{
				userShoppingListEntryRepository: NewDemoRepository(),
				userShoppingListRepository:      userShoppingList.NewDemoRepository(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDefaultController(
				tt.args.userShoppingListEntryRepository,
				tt.args.userShoppingListRepository,
			); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDefaultController() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultController_GetEntries(t *testing.T) {
	type fields struct {
		userShoppingListEntryRepository Repository
		userShoppingListRepository      userShoppingList.Repository
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
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"GET",
						"/api/v1/shoppinglistentries/1",
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
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"GET",
						"/api/v1/shoppinglistentries/1",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "1")
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
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"GET",
						"/api/v1/shoppinglistentries/1",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "1")
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
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"GET",
						"/api/v1/shoppinglistentries/1",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "1")
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
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"GET",
						"/api/v1/shoppinglistentries/abc",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "abc")
					ctx = context.WithValue(ctx, "auth_userId", uint64(2))
					ctx = context.WithValue(ctx, "auth_userRole", int64(1))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Unknown entries (expect 404)",
			fields: fields{
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"GET",
						"/api/v1/shoppinglistentries/4",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "4")
					ctx = context.WithValue(ctx, "auth_userId", uint64(2))
					ctx = context.WithValue(ctx, "auth_userRole", int64(1))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := DefaultController{
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			}
			controller.GetEntries(tt.args.writer, tt.args.request)
			if tt.args.writer.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, tt.args.writer.Code)
			}
		})
	}

	t.Run("Successfully get user shopping lists (expect 200 and shopping lists)", func(t *testing.T) {
		writer := httptest.NewRecorder()
		var request = httptest.NewRequest(
			"GET",
			"/api/v1/shoppinglistentries/1",
			nil)
		ctx := context.WithValue(request.Context(), "listId", "1")
		ctx = context.WithValue(ctx, "auth_userId", uint64(2))
		ctx = context.WithValue(ctx, "auth_userRole", int64(0))
		request = request.WithContext(ctx)

		controller := DefaultController{
			userShoppingListEntryRepository: setupMockEntryRepository(),
			userShoppingListRepository:      setupMockListRepository(),
		}

		controller.GetEntries(writer, request)

		if writer.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, writer.Code)
		}

		if writer.Header().Get("Content-Type") != "application/json" {
			t.Errorf("Expected content type %s, got %s",
				"application/json", writer.Header().Get("Content-Type"))
		}

		result := writer.Result()
		var response []entryModel.UserShoppingListEntry
		err := json.NewDecoder(result.Body).Decode(&response)
		if err != nil {
			t.Fatal(err.Error())
		}

		if len(response) != 2 {
			t.Errorf("Expected count of shopping lists is %d, got %d",
				2, len(response))
		}

		for i, entry := range response {
			if entry.ShoppingListId != response[i].ShoppingListId {
				t.Errorf("Expected userId of shopping lists is %d, got %d", entry.ShoppingListId, response[i].ShoppingListId)
			}
		}
	})
}

func TestDefaultController_GetEntry(t *testing.T) {
	type fields struct {
		userShoppingListEntryRepository Repository
		userShoppingListRepository      userShoppingList.Repository
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
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"GET",
						"/api/v1/shoppinglistentries/1/2",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "productId", "2")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Unauthorized, not matching userIds (expect 401)",
			fields: fields{
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"GET",
						"/api/v1/shoppinglistentries/1/2",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "productId", "2")
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
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"GET",
						"/api/v1/shoppinglistentries/1/2",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "productId", "2")
					ctx = context.WithValue(ctx, "auth_userId", uint64(2))
					ctx = context.WithValue(ctx, "auth_userRole", int64(0))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Valid request as admin(expect 200)",
			fields: fields{
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"GET",
						"/api/v1/shoppinglistentries/1/2",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "productId", "2")
					ctx = context.WithValue(ctx, "auth_userId", uint64(1))
					ctx = context.WithValue(ctx, "auth_userRole", int64(2))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid request, non-numeric listId (expect 400)",
			fields: fields{
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"GET",
						"/api/v1/shoppinglistentries/abc/abc",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "abc")
					ctx = context.WithValue(ctx, "productId", "abc")
					ctx = context.WithValue(ctx, "auth_userId", uint64(1))
					ctx = context.WithValue(ctx, "auth_userRole", int64(2))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Unknown entry (expect 404)",
			fields: fields{
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"GET",
						"/api/v1/shoppinglistentries/1/10",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "productId", "10")
					ctx = context.WithValue(ctx, "auth_userId", uint64(2))
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
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			}
			controller.GetEntry(tt.args.writer, tt.args.request)
			if tt.args.writer.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, tt.args.writer.Code)
			}
		})
	}
}

func TestDefaultController_PostEntry(t *testing.T) {
	type fields struct {
		userShoppingListEntryRepository Repository
		userShoppingListRepository      userShoppingList.Repository
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
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"POST",
						"/api/v1/shoppinglistentries/1/2",
						strings.NewReader(`{"count": 2, "note": "Test entry", "checked": false}`))
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "productId", "2")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Unauthorized, not matching userIds (expect 401)",
			fields: fields{
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"POST",
						"/api/v1/shoppinglistentries/1/3",
						strings.NewReader(`{"count": 2, "note": "Test entry", "checked": false}`))
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "productId", "3")
					ctx = context.WithValue(ctx, "auth_userId", uint64(1))
					ctx = context.WithValue(ctx, "auth_userRole", int64(0))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Valid request (expect 201)",
			fields: fields{
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"POST",
						"/api/v1/shoppinglistentries/1/3",
						strings.NewReader(`{"count": 2, "note": "Test entry", "checked": false}`))
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "productId", "3")
					ctx = context.WithValue(ctx, "auth_userId", uint64(2))
					ctx = context.WithValue(ctx, "auth_userRole", int64(1))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Valid request as admin(expect 201)",
			fields: fields{
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"POST",
						"/api/v1/shoppinglistentries/1/4",
						strings.NewReader(`{"count": 2, "note": "Test entry", "checked": false}`))
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "productId", "4")
					ctx = context.WithValue(ctx, "auth_userId", uint64(1))
					ctx = context.WithValue(ctx, "auth_userRole", int64(2))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid request, non-numeric listId (expect 400)",
			fields: fields{
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"POST",
						"/api/v1/shoppinglistentries/abc/abc",
						strings.NewReader(`{"count": 2, "note": "Test entry", "checked": false}`))
					ctx := context.WithValue(request.Context(), "listId", "abc")
					ctx = context.WithValue(ctx, "productId", "abc")
					ctx = context.WithValue(ctx, "auth_userId", uint64(1))
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
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			}
			controller.PostEntry(tt.args.writer, tt.args.request)
			if tt.args.writer.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, tt.args.writer.Code)
			}
		})
	}
}

func TestDefaultController_PutEntry(t *testing.T) {
	type fields struct {
		userShoppingListEntryRepository Repository
		userShoppingListRepository      userShoppingList.Repository
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
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/shoppinglistentries/1/2",
						strings.NewReader(`{"count": 10, "note": "Update entry", "checked": false}`))
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "productId", "2")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Unauthorized, not matching userIds (expect 401)",
			fields: fields{
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/shoppinglistentries/1/2",
						strings.NewReader(`{"count": 10, "note": "Update entry", "checked": false}`))
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "productId", "2")
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
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/shoppinglistentries/1/2",
						strings.NewReader(`{"count": 10, "note": "Update entry", "checked": false}`))
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "productId", "2")
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
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/shoppinglistentries/1/2",
						strings.NewReader(`{"count": 10, "note": "Update entry", "checked": false}`))
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "productId", "2")
					ctx = context.WithValue(ctx, "auth_userId", uint64(1))
					ctx = context.WithValue(ctx, "auth_userRole", int64(2))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid request, non-numeric listId (expect 400)",
			fields: fields{
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/shoppinglistentries/abc/abc",
						strings.NewReader(`{"count": 2, "note": "Test entry", "checked": false}`))
					ctx := context.WithValue(request.Context(), "listId", "abc")
					ctx = context.WithValue(ctx, "productId", "abc")
					ctx = context.WithValue(ctx, "auth_userId", uint64(1))
					ctx = context.WithValue(ctx, "auth_userRole", int64(2))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Unknown entry (expect 404)",
			fields: fields{
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"PUT",
						"/api/v1/shoppinglistentries/10/10",
						strings.NewReader(`{"count": 2, "note": "Test entry", "checked": false}`))
					ctx := context.WithValue(request.Context(), "listId", "10")
					ctx = context.WithValue(ctx, "productId", "10")
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
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			}
			controller.PutEntry(tt.args.writer, tt.args.request)
			if tt.args.writer.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, tt.args.writer.Code)
			}
		})
	}
}

func TestDefaultController_DeleteEntry(t *testing.T) {
	type fields struct {
		userShoppingListEntryRepository Repository
		userShoppingListRepository      userShoppingList.Repository
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
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"DELETE",
						"/api/v1/shoppinglistentries/1/2",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "productId", "2")
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Unauthorized, not matching userIds (expect 401)",
			fields: fields{
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"DELETE",
						"/api/v1/shoppinglistentries/1/2",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "productId", "2")
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
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"DELETE",
						"/api/v1/shoppinglistentries/1/2",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "productId", "2")
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
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"DELETE",
						"/api/v1/shoppinglistentries/1/2",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "1")
					ctx = context.WithValue(ctx, "productId", "2")
					ctx = context.WithValue(ctx, "auth_userId", uint64(1))
					ctx = context.WithValue(ctx, "auth_userRole", int64(2))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid request, non-numeric listId (expect 400)",
			fields: fields{
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"DELETE",
						"/api/v1/shoppinglistentries/1/2",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "abc")
					ctx = context.WithValue(ctx, "productId", "abc")
					ctx = context.WithValue(ctx, "auth_userId", uint64(1))
					ctx = context.WithValue(ctx, "auth_userRole", int64(2))
					return request.WithContext(ctx)
				}(),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Unknown entry (expect 404)",
			fields: fields{
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					var request = httptest.NewRequest(
						"DELETE",
						"/api/v1/shoppinglistentries/1/2",
						nil)
					ctx := context.WithValue(request.Context(), "listId", "10")
					ctx = context.WithValue(ctx, "productId", "10")
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
				userShoppingListEntryRepository: setupMockEntryRepository(),
				userShoppingListRepository:      setupMockListRepository(),
			}
			controller.DeleteEntry(tt.args.writer, tt.args.request)
			if tt.args.writer.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, tt.args.writer.Code)
			}
		})
	}
}

func setupMockEntryRepository() Repository {
	repository := NewDemoRepository()
	entries := setupMockEntrySlice()
	for _, entry := range entries {
		repository.Create(entry)
	}
	return repository
}

func setupMockEntrySlice() []*entryModel.UserShoppingListEntry {
	return []*entryModel.UserShoppingListEntry{
		{ShoppingListId: 1, ProductId: 1, Count: 3, Note: "Sample entry 1", Checked: false},
		{ShoppingListId: 1, ProductId: 2, Count: 2, Note: "Sample entry 2", Checked: true},
		{ShoppingListId: 2, ProductId: 1, Count: 1, Note: "Sample entry 3", Checked: false},
	}
}

func setupMockListRepository() userShoppingList.Repository {
	repository := userShoppingList.NewDemoRepository()
	lists := setupMockListSlice()
	for _, list := range lists {
		repository.Create(list)
	}
	return repository
}

func setupMockListSlice() []*listModel.UserShoppingList {
	return []*listModel.UserShoppingList{
		{Id: 1, UserId: 2, Description: "Frühstück mit Anne", Completed: false},
		{Id: 2, UserId: 2, Description: "Suppe", Completed: false},
		{Id: 3, UserId: 4, Description: "Geburtstagskuchen", Completed: false},
	}
}
