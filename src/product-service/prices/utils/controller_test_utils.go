package utils

import (
	"context"
	"net/http"
	"net/http/httptest"
)

func CreatePriceRequestWithValues(method string, path string, productId string, userId string) *http.Request {
	request := httptest.NewRequest(method, path, nil)
	ctx := request.Context()
	ctx = context.WithValue(ctx, "productId", productId)
	ctx = context.WithValue(ctx, "userId", userId)
	request = request.WithContext(ctx)
	return request
}
