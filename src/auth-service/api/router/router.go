package router

import (
	"net/http"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter(loginHandler http.Handler, registerHandler http.Handler) *Router {
	mux := http.NewServeMux()
	mux.Handle("/login", loginHandler)
	mux.Handle("/register", registerHandler)

	return &Router{
		mux: mux,
	}
}

func (handler *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.mux.ServeHTTP(w, r)
}
