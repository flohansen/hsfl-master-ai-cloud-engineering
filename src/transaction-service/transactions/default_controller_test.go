package transactions

import (
	"context"
	"encoding/json"
	"errors"
	mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/_mocks"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/transactions/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTransactionDefaultController(t *testing.T) {
	ctrl := gomock.NewController(t)
	transactionRepository := mocks.NewMockTransactionRepository(ctrl)
	controller := DefaultController{transactionRepository}

	t.Run("GetTransactions", func(t *testing.T) {
		t.Run("should return 500 INTERNAL SERVER ERROR if query failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/transactions", nil)

			transactionRepository.
				EXPECT().
				FindAll().
				Return(nil, errors.New("query failed")).
				Times(1)

			// when
			controller.GetTransactions(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return all transactions", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/transactions", nil)

			transactionRepository.
				EXPECT().
				FindAll().
				Return([]*model.Transaction{{ID: 999}}, nil).
				Times(1)

			// when
			controller.GetTransactions(w, r)

			// then
			res := w.Result()
			var response []model.Transaction
			err := json.NewDecoder(res.Body).Decode(&response)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Len(t, response, 1)
			assert.Equal(t, uint64(999), response[0].ID)
		})
	})

	t.Run("PostTransactions", func(t *testing.T) {
		t.Run("should return 400 BAD REQUEST if payload is not json", func(t *testing.T) {
			tests := []io.Reader{
				nil,
				strings.NewReader(`{"invalid`),
			}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/api/v1/transactions", test)

				// when
				controller.PostTransactions(w, r)

				// then
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})

		t.Run("should return 400 BAD REQUEST if payload is incomplete", func(t *testing.T) {
			tests := []io.Reader{
				strings.NewReader(`{"chapterId": 1}`),
			}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/api/v1/transactions", test)

				// when
				controller.PostTransactions(w, r)

				// then
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if persisting failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/transactions",
				strings.NewReader(`{"chapterID":1,"payingUserID":1, "receivingUserID":1,"amount":100}`))

			transactionRepository.
				EXPECT().
				Create([]*model.Transaction{{ChapterID: 1, PayingUserID: 1, ReceivingUserID: 1, Amount: 100}}).
				Return(errors.New("database error"))

			// when
			controller.PostTransactions(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should create new transaction", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/transactions",
				strings.NewReader(`{"chapterID":1,"payingUserID":1, "receivingUserID":1,"amount":100}`))

			transactionRepository.
				EXPECT().
				Create([]*model.Transaction{{ChapterID: 1, PayingUserID: 1, ReceivingUserID: 1, Amount: 100}}).
				Return(nil)

			// when
			controller.PostTransactions(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("GetTransaction", func(t *testing.T) {
		t.Run("should return 400 BAD REQUEST if transaction id is not numerical", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/transactions/aaa", nil)
			r = r.WithContext(context.WithValue(r.Context(), "transactionid", "aaa"))

			// when
			controller.GetTransaction(w, r)

			// then
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		t.Run("should return 500 INTERNAL SERVER ERROR query failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/transactions/1", nil)
			r = r.WithContext(context.WithValue(r.Context(), "transactionid", "1"))

			transactionRepository.
				EXPECT().
				FindById(uint64(1)).
				Return(nil, errors.New("database error"))

			// when
			controller.GetTransaction(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return 200 OK and transaction", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/transaction/1", nil)
			r = r.WithContext(context.WithValue(r.Context(), "transactionid", "1"))

			transactionRepository.
				EXPECT().
				FindById(uint64(1)).
				Return(&model.Transaction{ID: 1, ChapterID: 1, PayingUserID: 1}, nil)

			// when
			controller.GetTransaction(w, r)

			// then
			res := w.Result()
			var response model.Transaction
			err := json.NewDecoder(res.Body).Decode(&response)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, uint64(1), response.ID)
			assert.Equal(t, uint64(1), response.ChapterID)
		})
	})
}
