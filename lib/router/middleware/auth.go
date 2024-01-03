package middleware

import (
	router "hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib"
	"net/http"
)

func CreateAuthMiddleware() router.Middleware {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		return r
	}
}
