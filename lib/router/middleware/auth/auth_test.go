package auth

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/rpc/user"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/rpc/user/_mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateAuthMiddleware(t *testing.T) {
	userServiceClient := &_mocks.MockUserServiceClient{}
	authMiddleware := CreateAuthMiddleware(userServiceClient)

	validToken := "valid-token"
	invalidToken := "invalid-token"

	expectedUser := &user.User{
		Id:    1,
		Email: "ada.lovelace@gmail.com",
		Name:  "Ada Lovelace",
		Role:  0,
	}

	userServiceClient.On("ValidateUserToken", mock.Anything,
		&user.ValidateUserTokenRequest{
			Token: validToken}).
		Return(&user.ValidateUserTokenResponse{
			User: expectedUser,
		}, nil)

	userServiceClient.On("ValidateUserToken", mock.Anything,
		&user.ValidateUserTokenRequest{
			Token: invalidToken}).
		Return(nil, status.Error(codes.Unauthenticated, "Verification failed"))

	t.Run("ValidToken", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/secret", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)

		updatedReq := authMiddleware(nil, req)

		_, userIDExists := updatedReq.Context().Value("auth_userId").(uint64)
		_, userEmailExists := updatedReq.Context().Value("auth_userEmail").(string)
		_, userNameExists := updatedReq.Context().Value("auth_userName").(string)
		_, userRoleExists := updatedReq.Context().Value("auth_userRole").(int64)

		assert.True(t, userIDExists, "User ID should exist")
		assert.True(t, userEmailExists, "User email should exist")
		assert.True(t, userNameExists, "User name should exist")
		assert.True(t, userRoleExists, "User role should exist")
	})

	t.Run("InvalidToken", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/secret", nil)
		req.Header.Set("Authorization", "Bearer "+invalidToken)

		updatedReq := authMiddleware(nil, req)

		_, userIDExists := updatedReq.Context().Value("auth_userId").(uint64)
		_, userEmailExists := updatedReq.Context().Value("auth_userEmail").(string)
		_, userNameExists := updatedReq.Context().Value("auth_userName").(string)
		_, userRoleExists := updatedReq.Context().Value("auth_userRole").(int64)

		assert.False(t, userIDExists, "User ID should not exist")
		assert.False(t, userEmailExists, "User email should not exist")
		assert.False(t, userNameExists, "User name should not exist")
		assert.False(t, userRoleExists, "User role should not exist")
	})

	userServiceClient.AssertExpectations(t)
}

func Test_getToken(t *testing.T) {
	tests := []struct {
		name          string
		authHeader    string
		expectedToken string
		expectError   bool
	}{
		{
			name:          "Valid Authorization Header",
			authHeader:    "Bearer validtoken123",
			expectedToken: "validtoken123",
			expectError:   false,
		},
		{
			name:        "Empty Authorization Header",
			authHeader:  "",
			expectError: true,
		},
		{
			name:        "Invalid Format Authorization Header",
			authHeader:  "Invalidformat",
			expectError: true,
		},
		{
			name:        "Incorrect Scheme",
			authHeader:  "Basic validtoken123",
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if test.authHeader != "" {
				req.Header.Add("Authorization", test.authHeader)
			}

			token, err := getToken(req)
			if test.expectError {
				if err == nil {
					t.Errorf("Expected an error but didn't get one")
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error but got one: %v", err)
				}
				if token != test.expectedToken {
					t.Errorf("Expected token %s, got %s", test.expectedToken, token)
				}
			}
		})
	}
}
