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

func TestChapterDefaultController(t *testing.T) {
	ctrl := gomock.NewController(t)

	chapterRepository := mocks.NewMockChapterRepository(ctrl)
	controller := DefaultController{nil, chapterRepository}

	t.Run("GetChapters", func(t *testing.T) {
		t.Run("should return 500 INTERNAL SERVER ERROR if query failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/books/1/chapters", nil)

			chapterRepository.
				EXPECT().
				FindAll().
				Return(nil, errors.New("query failed")).
				Times(1)

			// when
			controller.GetChapters(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return all chapters", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/books/1/chapters", nil)

			chapterRepository.
				EXPECT().
				FindAll().
				Return([]*model.Chapter{{ID: 999}}, nil).
				Times(1)

			// when
			controller.GetChapters(w, r)

			// then
			res := w.Result()
			var response []model.Chapter
			err := json.NewDecoder(res.Body).Decode(&response)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Len(t, response, 1)
			assert.Equal(t, int64(999), response[0].ID)
		})
	})

	t.Run("PostChapters", func(t *testing.T) {
		t.Run("should return 400 BAD REQUEST if payload is not json", func(t *testing.T) {
			tests := []io.Reader{
				nil,
				strings.NewReader(`{"invalid`),
			}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/api/v1/books/1/chapters", test)

				// when
				controller.PostChapters(w, r)

				// then
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})

		t.Run("should return 400 BAD REQUEST if payload is incomplete", func(t *testing.T) {
			tests := []io.Reader{
				strings.NewReader(`{"price": 99.99}`),
				strings.NewReader(`{"description": "amazing chapter"}`),
				strings.NewReader(`{"author": "the best author"}`),
			}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/api/v1/books/1/chapters", test)

				// when
				controller.PostChapters(w, r)

				// then
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if persisting failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/books/1/chapters",
				strings.NewReader(`{"name":"test chapter","author":"the best author"}`))

			chapterRepository.
				EXPECT().
				Create([]*model.Chapter{{Name: "test chapter", Author: "the best author"}}).
				Return(errors.New("database error"))

			// when
			controller.PostChapters(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should create new chapter", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/books/1/chapters",
				strings.NewReader(`{"name":"test chapter","author":"the best author"}`))

			chapterRepository.
				EXPECT().
				Create([]*model.Chapter{{Name: "test chapter", Author: "the best author"}}).
				Return(nil)

			// when
			controller.PostChapters(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("GetChapter", func(t *testing.T) {
		t.Run("should return 400 BAD REQUEST if chapter id is not numerical", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/books/1/chapters/aaa", nil)
			r = r.WithContext(context.WithValue(r.Context(), "chapterid", "aaa"))

			// when
			controller.GetChapter(w, r)

			// then
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		t.Run("should return 500 INTERNAL SERVER ERROR query failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/books/1/chapters/1", nil)
			r = r.WithContext(context.WithValue(r.Context(), "chapterid", "1"))

			chapterRepository.
				EXPECT().
				FindById(int64(1)).
				Return(nil, errors.New("database error"))

			// when
			controller.GetChapter(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return 200 OK and chapter", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/books/1/chapters/1", nil)
			r = r.WithContext(context.WithValue(r.Context(), "chapterid", "1"))

			chapterRepository.
				EXPECT().
				FindById(int64(1)).
				Return(&model.Chapter{ID: 1, Name: "test chapter", Author: "the best author"}, nil)

			// when
			controller.GetChapter(w, r)

			// then
			res := w.Result()
			var response model.Chapter
			err := json.NewDecoder(res.Body).Decode(&response)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, int64(1), response.ID)
			assert.Equal(t, "test chapter", response.Name)
		})
	})

	t.Run("PutChapter", func(t *testing.T) {
		t.Run("should return 400 BAD REQUEST if chapter id is not numerical", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/books/1/chapters/aaa", nil)
			r = r.WithContext(context.WithValue(r.Context(), "chapterid", "aaa"))

			// when
			controller.PutChapter(w, r)

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
				r := httptest.NewRequest("PUT", "/api/v1/books/1/chapters/1", test)
				r = r.WithContext(context.WithValue(r.Context(), "chapterid", "1"))

				// when
				controller.PutChapter(w, r)

				// then
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if query failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/books/1/chapters/1",
				strings.NewReader(`{"id": 999}`))
			r = r.WithContext(context.WithValue(r.Context(), "chapterid", "1"))

			chapterRepository.
				EXPECT().
				Create([]*model.Chapter{{ID: 1}}).
				Return(errors.New("database error"))

			// when
			controller.PutChapter(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should update one chapter", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/books/1/chapters/1",
				strings.NewReader(`{"id": 999}`))
			r = r.WithContext(context.WithValue(r.Context(), "chapterid", "1"))

			chapterRepository.
				EXPECT().
				Create([]*model.Chapter{{ID: 1}}).
				Return(nil)

			// when
			controller.PutChapter(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("DeleteChapter", func(t *testing.T) {
		t.Run("should return 400 BAD REQUEST if chapter id is not numerical", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/books/1/chapters/aaa", nil)
			r = r.WithContext(context.WithValue(r.Context(), "chapterid", "aaa"))

			// when
			controller.DeleteChapter(w, r)

			// then
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if query fails", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/books/1/chapters/1", nil)
			r = r.WithContext(context.WithValue(r.Context(), "chapterid", "1"))

			chapterRepository.
				EXPECT().
				Delete([]*model.Chapter{{ID: 1}}).
				Return(errors.New("database error"))

			// when
			controller.DeleteChapter(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return 200 OK", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/books/1/chapters/1", nil)
			r = r.WithContext(context.WithValue(r.Context(), "chapterid", "1"))

			chapterRepository.
				EXPECT().
				Delete([]*model.Chapter{{ID: 1}}).
				Return(nil)

			// when
			controller.DeleteChapter(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
}
