package router

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
)

type Middleware = func(w http.ResponseWriter, r *http.Request) *http.Request

type route struct {
	method  string
	pattern *regexp.Regexp
	handler http.HandlerFunc
	params  []string
}

type Router struct {
	routes      []route
	middlewares []Middleware
}

func New() *Router {
	return &Router{}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range router.routes {
		if r.Method != route.method {
			continue
		}

		matches := route.pattern.FindStringSubmatch(r.URL.Path)

		if len(matches) > 0 {
			r = createRequestContext(r, route.params, matches[1:])
			for _, middleware := range router.middlewares {
				r = middleware(w, r)
			}
			route.handler(w, r)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func createRequestContext(r *http.Request, paramKeys []string, paramValues []string) *http.Request {
	if len(paramKeys) == 0 {
		return r
	}

	ctx := r.Context()
	for i := 0; i < len(paramKeys); i++ {
		ctx = context.WithValue(ctx, paramKeys[i], paramValues[i])
	}

	return r.WithContext(ctx)
}

func (router *Router) addRoute(method string, pattern string, handler http.HandlerFunc, middlewares ...Middleware) {
	paramMatcher := regexp.MustCompile(":([a-zA-Z]+)|\\*")
	paramMatches := paramMatcher.FindAllStringSubmatch(pattern, -1)

	params := make([]string, len(paramMatches))

	if len(paramMatches) > 0 {
		pattern = paramMatcher.ReplaceAllStringFunc(pattern, func(substring string) string {
			if substring == "*" {
				return "(.*)"
			} else {
				return "([^/]+)/?"
			}
		})

		wildcardNo := 0
		for i, match := range paramMatches {
			if match[1] != "" {
				params[i] = match[1]
			} else {
				params[i] = fmt.Sprintf("wildcard%d", wildcardNo)
				wildcardNo++
			}

		}
	}

	router.routes = append(router.routes, route{
		method:  method,
		pattern: regexp.MustCompile("^" + pattern + "$"),
		handler: func(w http.ResponseWriter, r *http.Request) {
			for _, middleware := range middlewares {
				r = middleware(w, r)
			}
			handler(w, r)
		},
		params: params,
	})
}

func (router *Router) GET(pattern string, handler http.HandlerFunc, middlewares ...Middleware) {
	router.addRoute(http.MethodGet, pattern, handler, middlewares...)
}

func (router *Router) POST(pattern string, handler http.HandlerFunc, middlewares ...Middleware) {
	router.addRoute(http.MethodPost, pattern, handler, middlewares...)
}

func (router *Router) PUT(pattern string, handler http.HandlerFunc, middlewares ...Middleware) {
	router.addRoute(http.MethodPut, pattern, handler, middlewares...)
}

func (router *Router) DELETE(pattern string, handler http.HandlerFunc, middlewares ...Middleware) {
	router.addRoute(http.MethodDelete, pattern, handler, middlewares...)
}

func (router *Router) PATCH(pattern string, handler http.HandlerFunc, middlewares ...Middleware) {
	router.addRoute(http.MethodPatch, pattern, handler, middlewares...)
}

func (router *Router) CONNECT(pattern string, handler http.HandlerFunc, middlewares ...Middleware) {
	router.addRoute(http.MethodConnect, pattern, handler, middlewares...)
}

func (router *Router) HEAD(pattern string, handler http.HandlerFunc, middlewares ...Middleware) {
	router.addRoute(http.MethodHead, pattern, handler, middlewares...)
}

func (router *Router) OPTIONS(pattern string, handler http.HandlerFunc, middlewares ...Middleware) {
	router.addRoute(http.MethodOptions, pattern, handler, middlewares...)
}

func (router *Router) TRACE(pattern string, handler http.HandlerFunc, middlewares ...Middleware) {
	router.addRoute(http.MethodTrace, pattern, handler, middlewares...)
}

func (router *Router) ALL(pattern string, handler http.HandlerFunc, middlewares ...Middleware) {
	router.GET(pattern, handler, middlewares...)
	router.POST(pattern, handler, middlewares...)
	router.PUT(pattern, handler, middlewares...)
	router.DELETE(pattern, handler, middlewares...)
	router.PATCH(pattern, handler, middlewares...)
	router.CONNECT(pattern, handler, middlewares...)
	router.HEAD(pattern, handler, middlewares...)
	router.OPTIONS(pattern, handler, middlewares...)
	router.TRACE(pattern, handler, middlewares...)
}

func (router *Router) RegisterMiddleware(middleware Middleware) {
	router.middlewares = append(router.middlewares, middleware)
}
