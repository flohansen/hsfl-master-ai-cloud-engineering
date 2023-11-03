package transactions

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/transactions/model"
	"github.com/stretchr/testify/assert"
)

func TestPsqlRepository(t *testing.T) {
	db, dbmock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	repository := PsqlRepository{db}

	t.Run("Create", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			transactions := []*model.Transaction{{
				ID:           1,
				ChapterID:    1,
				PayingUserID: 1,
				Amount:       0,
			}}

			dbmock.
				ExpectExec(`insert into transactions`).
				WillReturnError(errors.New("database error"))

			// when
			err := repository.Create(transactions)

			// then
			assert.Error(t, err)
		})

		t.Run("should insert transactions in batches", func(t *testing.T) {
			// given
			transactions := []*model.Transaction{
				{
					ID:           1,
					ChapterID:    1,
					PayingUserID: 1,
					Amount:       0,
				},
				{
					ID:           2,
					ChapterID:    2,
					PayingUserID: 1,
					Amount:       100,
				},
			}

			dbmock.
				ExpectExec(`insert into transactions \(chapterid, payinguserid, amount\) values \(\$1,\$2,\$3\),\(\$4,\$5,\$6\)`).
				WithArgs(1, 1, 0, 2, 1, 100).
				WillReturnResult(sqlmock.NewResult(0, 2))

			// when
			err := repository.Create(transactions)

			// then
			assert.NoError(t, err)
		})
	})

	t.Run("FindAll", func(t *testing.T) {
		t.Run("should return all transactions", func(t *testing.T) {
			// given
			dbmock.ExpectQuery(`select id, chapterid, payinguserid, amount from transactions`).
				WillReturnRows(sqlmock.NewRows([]string{"id", "chapterid", "payinguserid", "amount"}).
					AddRow(1, 1, 1, 0).
					AddRow(2, 2, 2, 1))

			// when
			transactions, err := repository.FindAll()

			// then
			assert.NoError(t, err)
			assert.NoError(t, dbmock.ExpectationsWereMet())
			assert.Len(t, transactions, 2)
			assert.Equal(t, uint64(1), transactions[0].ID)
			assert.Equal(t, uint64(2), transactions[1].ID)
		})
	})

	t.Run("FindByID", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			id := uint64(1)

			dbmock.
				ExpectQuery(`select id, chapterid, payinguserid, amount from transactions where id = \$1`).
				WillReturnError(errors.New("database error"))

			// when
			transaction, err := repository.FindById(id)

			// then
			assert.Error(t, err)
			assert.Nil(t, transaction)
		})

		t.Run("should return transactions by id", func(t *testing.T) {
			// given
			id := uint64(1)

			dbmock.
				ExpectQuery(`select id, chapterid, payinguserid, amount from transactions where id = \$1`).
				WillReturnRows(sqlmock.NewRows([]string{"id", "chapterid", "payinguserid", "amount"}).
					AddRow(1, 1, 1, 0))

			// when
			transaction, err := repository.FindById(id)

			// then
			assert.NoError(t, err)
			assert.Equal(t, transaction.ID, id)
		})
	})
}
