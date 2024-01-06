package test

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router"
	"net/http"
)

func CreateEmptyMiddleware() router.Middleware {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		return r
	}
}
