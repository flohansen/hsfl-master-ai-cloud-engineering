package router

import (
	"context"
	"fmt"
	mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	ctrl := gomock.NewController(t)

	transactionsController := mocks.NewMockController(ctrl)
	router := New(transactionsController)

	t.Run("/api/v1/transactions", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not GET or POST", func(t *testing.T) {
			tests := []string{"DELETE", "PUT", "HEAD", "CONNECT", "OPTIONS", "TRACE", "PATCH"}
			for _, test := range tests {
				fmt.Println(test)
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/api/v1/books", nil)

				// when
				router.ServeHTTP(w, r)

				// then
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		})
		t.Run("should call GET handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/transactions", nil)

			fmt.Println("transactionsController")
			fmt.Println(w, r)
			transactionsController.
				EXPECT().
				GetTransactions(w, r).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call POST handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/transactions", nil)

			transactionsController.
				EXPECT().
				PostTransactions(w, r).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/api/v1/transactions/:transactionid", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not GET", func(t *testing.T) {
			tests := []string{"POST", "HEAD", "CONNECT", "OPTIONS", "TRACE", "PATCH", "PUT", "DELETE"}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/api/v1/transactions/1", nil)

				// when
				router.ServeHTTP(w, r)

				// then
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		})

		t.Run("should call GET handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/transactions/1", nil)

			transactionsController.
				EXPECT().
				GetTransaction(w, r.WithContext(context.WithValue(r.Context(), "transactionid", "1"))).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
}
