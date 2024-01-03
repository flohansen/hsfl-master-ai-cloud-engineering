package auth

import (
	"context"
	"fmt"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/rpc/user"
	"net/http"
	"strings"
)

func CreateAuthMiddleware(grpcUserServiceClient user.UserServiceClient) router.Middleware {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		token, err := getToken(r)
		if err != nil {
			println(err)
		} else {
			usr, err := grpcUserServiceClient.ValidateUserToken(context.Background(), &user.ValidateUserTokenRequest{
				Token: token,
			})
			if err != nil {
				println(err)
			} else {
				ctx := context.WithValue(r.Context(), "auth_userId", usr.User.Id)
				ctx = context.WithValue(ctx, "auth_userEmail", usr.User.Email)
				ctx = context.WithValue(ctx, "auth_userName", usr.User.Name)
				ctx = context.WithValue(ctx, "auth_userRole", usr.User.Role)
				return r.WithContext(ctx)
			}
		}
		return r
	}
}

func getToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is not provided")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("authorization header format is not valid")
	}

	return parts[1], nil
}
