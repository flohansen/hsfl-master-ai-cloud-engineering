package books

import (
	"context"
	"encoding/json"
	"errors"
	mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/_mocks"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBookDefaultController(t *testing.T) {
	ctrl := gomock.NewController(t)
	bookRepository := mocks.NewMockBookRepository(ctrl)
	controller := DefaultController{bookRepository, nil}

	t.Run("GetBooks", func(t *testing.T) {
		t.Run("should return 500 INTERNAL SERVER ERROR if query failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/books", nil)

			bookRepository.
				EXPECT().
				FindAll().
				Return(nil, errors.New("query failed")).
				Times(1)

			// when
			controller.GetBooks(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return all books", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/books", nil)

			bookRepository.
				EXPECT().
				FindAll().
				Return([]*model.Book{{ID: 999}}, nil).
				Times(1)

			// when
			controller.GetBooks(w, r)

			// then
			res := w.Result()
			var response []model.Book
			err := json.NewDecoder(res.Body).Decode(&response)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Len(t, response, 1)
			assert.Equal(t, uint64(999), response[0].ID)
		})
	})

	t.Run("PostBooks", func(t *testing.T) {
		t.Run("should return 400 BAD REQUEST if payload is not json", func(t *testing.T) {
			tests := []io.Reader{
				nil,
				strings.NewReader(`{"invalid`),
			}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/api/v1/books", test)

				// when
				controller.PostBooks(w, r)

				// then
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})

		t.Run("should return 400 BAD REQUEST if payload is incomplete", func(t *testing.T) {
			tests := []io.Reader{
				strings.NewReader(`{"description": "amazing book"}`),
				strings.NewReader(`{"authorid": 1}`),
			}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/api/v1/books", test)

				// when
				controller.PostBooks(w, r)

				// then
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if persisting failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/books",
				strings.NewReader(`{"name":"test book","authorid":1}`))

			bookRepository.
				EXPECT().
				Create([]*model.Book{{Name: "test book", AuthorID: 1}}).
				Return(errors.New("database error"))

			// when
			controller.PostBooks(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should create new book", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/books",
				strings.NewReader(`{"name":"test book","authorid":1}`))

			bookRepository.
				EXPECT().
				Create([]*model.Book{{Name: "test book", AuthorID: 1}}).
				Return(nil)

			// when
			controller.PostBooks(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("GetBook", func(t *testing.T) {
		t.Run("should return 400 BAD REQUEST if book id is not numerical", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/books/aaa", nil)
			r = r.WithContext(context.WithValue(r.Context(), "bookid", "aaa"))

			// when
			controller.GetBook(w, r)

			// then
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		t.Run("should return 500 INTERNAL SERVER ERROR query failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/books/1", nil)
			r = r.WithContext(context.WithValue(r.Context(), "bookid", "1"))

			bookRepository.
				EXPECT().
				FindById(uint64(1)).
				Return(nil, errors.New("database error"))

			// when
			controller.GetBook(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return 200 OK and book", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/books/1", nil)
			r = r.WithContext(context.WithValue(r.Context(), "bookid", "1"))

			bookRepository.
				EXPECT().
				FindById(uint64(1)).
				Return(&model.Book{ID: 1, Name: "test book", AuthorID: 1}, nil)

			// when
			controller.GetBook(w, r)

			// then
			res := w.Result()
			var response model.Book
			err := json.NewDecoder(res.Body).Decode(&response)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, uint64(1), response.ID)
			assert.Equal(t, uint64(1), response.AuthorID)
			assert.Equal(t, "test book", response.Name)
		})
	})

	t.Run("PutBook", func(t *testing.T) {
		t.Run("should return 400 BAD REQUEST if book id is not numerical", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/books/aaa", nil)
			r = r.WithContext(context.WithValue(r.Context(), "bookid", "aaa"))

			// when
			controller.PutBook(w, r)

			// then
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		t.Run("should return 400 BAD REQUEST if payload is not json", func(t *testing.T) {
			tests := []io.Reader{
				nil,
				strings.NewReader(`{"invalid`),
			}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("PUT", "/api/v1/books/1", test)
				r = r.WithContext(context.WithValue(r.Context(), "bookid", "1"))

				// when
				controller.PutBook(w, r)

				// then
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if query failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/books/1",
				strings.NewReader(`{"id": 999}`))
			r = r.WithContext(context.WithValue(r.Context(), "bookid", "1"))

			bookRepository.
				EXPECT().
				Update(uint64(1), &model.UpdateBook{}).
				Return(errors.New("database error"))

			// when
			controller.PutBook(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should update one book", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/books/1",
				strings.NewReader(`{"id": 999}`))
			r = r.WithContext(context.WithValue(r.Context(), "bookid", "1"))

			bookRepository.
				EXPECT().
				Update(uint64(1), &model.UpdateBook{}).
				Return(nil)

			// when
			controller.PutBook(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("DeleteBook", func(t *testing.T) {
		t.Run("should return 400 BAD REQUEST if book id is not numerical", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/books/aaa", nil)
			r = r.WithContext(context.WithValue(r.Context(), "bookid", "aaa"))

			// when
			controller.DeleteBook(w, r)

			// then
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if query fails", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/books/1", nil)
			r = r.WithContext(context.WithValue(r.Context(), "bookid", "1"))

			bookRepository.
				EXPECT().
				Delete([]*model.Book{{ID: 1}}).
				Return(errors.New("database error"))

			// when
			controller.DeleteBook(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return 200 OK", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/books/1", nil)
			r = r.WithContext(context.WithValue(r.Context(), "bookid", "1"))

			bookRepository.
				EXPECT().
				Delete([]*model.Book{{ID: 1}}).
				Return(nil)

			// when
			controller.DeleteBook(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
}
