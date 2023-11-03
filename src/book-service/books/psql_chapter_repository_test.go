package books

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books/model"
	"github.com/stretchr/testify/assert"
)

func TestPsqlChapterRepository(t *testing.T) {
	db, dbmock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	repository := PsqlChapterRepository{db}

	t.Run("Create", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			chapters := []*model.Chapter{{
				ID:      1,
				BookID:  1,
				Name:    "doesnt matter",
				Price:   0,
				Content: "doesnt matter",
			}}

			dbmock.
				ExpectExec(`insert into chapters`).
				WillReturnError(errors.New("database error"))

			// when
			err := repository.Create(chapters)

			// then
			assert.Error(t, err)
		})

		t.Run("should insert chapters in batches", func(t *testing.T) {
			// given
			chapters := []*model.Chapter{
				{
					ID:      1,
					BookID:  1,
					Name:    "doesnt matter",
					Price:   0,
					Content: "doesnt matter",
				},
				{
					ID:      2,
					BookID:  1,
					Name:    "doesnt matter",
					Price:   0,
					Content: "doesnt matter",
				},
			}

			dbmock.
				ExpectExec(`insert into chapters \(bookId, name, price, content\) values \(\$1,\$2,\$3,\$4\),\(\$5,\$6,\$7,\$8\)`).
				WithArgs(1, "doesnt matter", 0, "doesnt matter", 1, "doesnt matter", 0, "doesnt matter").
				WillReturnResult(sqlmock.NewResult(0, 2))

			// when
			err := repository.Create(chapters)

			// then
			assert.NoError(t, err)
		})
	})

	t.Run("FindByID", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			var id uint64 = 1

			dbmock.
				ExpectQuery(`select id, bookId, name, price, content from chapters where id = \$1`).
				WillReturnError(errors.New("database error"))

			// when
			chapters, err := repository.FindById(id)

			// then
			assert.Error(t, err)
			assert.Nil(t, chapters)
		})

		t.Run("should return chapters by id", func(t *testing.T) {
			// given
			var id uint64 = 1

			dbmock.
				ExpectQuery(`select id, bookId, name, price, content from chapters where id = \$1`).
				WillReturnRows(sqlmock.NewRows([]string{"id", "bookId", "name", "price", "content"}).
					AddRow(1, 1, "doesnt matter", 0, "doesnt matter"))

			// when
			chapter, err := repository.FindById(id)

			// then
			assert.NoError(t, err)
			assert.Equal(t, chapter.ID, id)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			chapters := []*model.Chapter{
				{
					ID:      1,
					BookID:  1,
					Name:    "doesnt matter",
					Price:   0,
					Content: "doesnt matter",
				},
				{
					ID:      2,
					BookID:  1,
					Name:    "doesnt matter",
					Price:   0,
					Content: "doesnt matter",
				},
			}

			dbmock.
				ExpectExec(`delete from chapters`).
				WillReturnError(errors.New("database error"))

			// when
			err := repository.Delete(chapters)

			// then
			assert.Error(t, err)
		})

		t.Run("should delete chapters in batches", func(t *testing.T) {
			// given
			chapters := []*model.Chapter{
				{
					ID:      1,
					BookID:  1,
					Name:    "doesnt matter",
					Price:   0,
					Content: "doesnt matter",
				},
				{
					ID:      2,
					BookID:  1,
					Name:    "doesnt matter",
					Price:   0,
					Content: "doesnt matter",
				},
			}

			dbmock.
				ExpectExec(`delete from chapters where id in \(\$1,\$2\)`).
				WithArgs(1, 2).
				WillReturnResult(sqlmock.NewResult(0, 2))

			// when
			err := repository.Delete(chapters)

			// then
			assert.NoError(t, err)
		})
	})
	t.Run("Update", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			newChapterData := &model.UpdateChapter{
				Name:    "Updated Chapter",
				Price:   100,
				Content: "This is a new text",
			}

			dbmock.
				ExpectExec(`update chapter set name = \$1, price = \$2, content = \$3 where id = \$4`).
				WillReturnError(errors.New("database error"))

			// when
			err := repository.Update(1, newChapterData)

			// then
			assert.Error(t, err)
		})

		t.Run("should update chapter", func(t *testing.T) {
			// given
			newChapterData := &model.UpdateChapter{
				Name:    "Updated Chapter",
				Price:   100,
				Content: "This is a new text",
			}

			dbmock.
				ExpectExec(`update chapter set name = \$1, price = \$2, content = \$3 where id = \$4`).
				WithArgs("Updated Chapter", 100, "This is a new text", 1).
				WillReturnResult(sqlmock.NewResult(0, 1))

			// when
			err := repository.Update(1, newChapterData)

			// then
			assert.Error(t, err)
		})
	})
}
