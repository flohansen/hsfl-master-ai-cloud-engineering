package middleware

import (
	"context"
	"fmt"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/auth"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func CreateLocalAuthMiddleware(userRepository *user.Repository, generator auth.TokenGenerator) router.Middleware {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		token, err := getToken(r)
		if err != nil {
			println(err)
		} else {
			claims, err := generator.VerifyToken(token)
			if err != nil {
				log.Println("Verification failed: ", err.Error())
				return r
			}

			id, err := strconv.ParseUint(claims["id"].(string), 10, 64)
			if err != nil {
				log.Println("Verification failed: Can't find user id in claim.")
				return r
			}

			user, err := (*userRepository).FindById(id)
			if err != nil {
				log.Println("Verification failed: ", err.Error())
				return r
			}
			if err != nil {
				println(err)
			} else {
				ctx := context.WithValue(r.Context(), "auth_userId", user.Id)
				ctx = context.WithValue(ctx, "auth_userEmail", user.Email)
				ctx = context.WithValue(ctx, "auth_userName", user.Name)
				ctx = context.WithValue(ctx, "auth_userRole", user.Role)
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
