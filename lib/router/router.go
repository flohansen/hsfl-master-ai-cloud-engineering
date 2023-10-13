package router

import (
	"context"
	"net/http"
	"regexp"
)

type route struct {
	method  string
	handler http.HandlerFunc
	pattern *regexp.Regexp
	params  []string
}

type router struct {
	routes []route
}

func New() *router {
	return &router{}
}

func createContext(r *http.Request, params []string, paramValues []string) *http.Request {
	ctx := r.Context()
	for i, param := range params {
		ctx = context.WithValue(ctx, param, paramValues[i])
	}

	return r.WithContext(ctx)
}

func (router *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range router.routes {
		if r.Method != route.method {
			continue
		}
		matches := route.pattern.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			r = createContext(r, route.params, matches[1:])
			route.handler(w, r)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func (router *router) addRoute(method string, target string, handler http.HandlerFunc) {
	patternMatcher := regexp.MustCompile(":([a-z]+)")
	matches := patternMatcher.FindAllStringSubmatch(target, -1)
	params := make([]string, len(matches))

	var targetPattern string
	if len(matches) > 0 {
		targetPattern = patternMatcher.ReplaceAllString(target, "([^/]+)")
		for i, match := range matches {
			params[i] = match[1]
		}
	} else {
		targetPattern = target
	}

	router.routes = append(router.routes, route{
		method:  method,
		handler: handler,
		pattern: regexp.MustCompile("^" + targetPattern + "$"),
		params:  params,
	})
}

func (router *router) GET(target string, handler http.HandlerFunc) {
	router.addRoute(http.MethodGet, target, handler)
}

func (router *router) POST(target string, handler http.HandlerFunc) {
	router.addRoute(http.MethodPost, target, handler)
}

func (router *router) DELETE(target string, handler http.HandlerFunc) {
	router.addRoute(http.MethodDelete, target, handler)
}

func (router *router) PUT(target string, handler http.HandlerFunc) {
	router.addRoute(http.MethodPut, target, handler)
}
