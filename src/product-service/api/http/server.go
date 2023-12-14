package http

import (
	"bytes"
	"fmt"
	"golang.org/x/sync/singleflight"
	"hash/fnv"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/api/http/router"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

type server struct {
	httpRouter *router.Router
	g          *singleflight.Group
}

func NewServer(productController *products.Controller, pricesController *prices.Controller) *server {
	httpRouter := router.New(productController, pricesController)
	g := &singleflight.Group{}
	return &server{httpRouter, g}
}

func (serv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hash, err := hashRequestFast(r)

	response, err, _ := serv.g.Do(hash, func() (interface{}, error) {
		rw := httptest.NewRecorder()
		serv.httpRouter.ServeHTTP(rw, r)
		return rw, nil
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for k, v := range response.(*httptest.ResponseRecorder).Result().Header {
		w.Header()[k] = v
	}
	w.WriteHeader(response.(*httptest.ResponseRecorder).Code)
	response.(*httptest.ResponseRecorder).Body.WriteTo(w)
}

func hashRequestFast(r *http.Request) (string, error) {
	var sb strings.Builder

	// Method and URL
	sb.WriteString(r.Method)
	sb.WriteString(r.URL.String())

	// Read and hash the body if it's not nil
	if r.Body != nil {
		// Read the body
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return "", err
		}

		// Hash the body
		_, err = sb.Write(bodyBytes)
		if err != nil {
			return "", err
		}

		// Restore the body for future use
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	// Compute the hash
	hasher := fnv.New32a()
	_, err := hasher.Write([]byte(sb.String()))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
