package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/lib/router"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/lib/rpc/auth"
)

type UserIdKey string

func CreateAuthMiddleware(authServiceClient auth.AuthServiceClient) router.MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			tokenParts := strings.Split(token, " ")

			if len(tokenParts) != 2 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			bearerToken := tokenParts[1]

			resp, err := authServiceClient.ValidateToken(context.Background(), &auth.ValidateTokenRequest{
				Token: bearerToken,
			})

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			var user_id_key UserIdKey = "user_id"

			ctx := r.Context()
			ctx = context.WithValue(ctx, user_id_key, resp.User.Id)
			r = r.WithContext(ctx)

			next(w, r)
		}
	}
}
