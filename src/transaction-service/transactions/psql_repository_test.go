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
				ID:              1,
				ChapterID:       1,
				PayingUserID:    1,
				ReceivingUserID: 1,
				Amount:          0,
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
					ID:              1,
					ChapterID:       1,
					PayingUserID:    1,
					ReceivingUserID: 1,
					Amount:          0,
				},
				{
					ID:              2,
					ChapterID:       2,
					PayingUserID:    1,
					ReceivingUserID: 1,
					Amount:          100,
				},
			}

			dbmock.
				ExpectExec(`insert into transactions \(chapterid, payinguserid, receivinguserid, amount\) values \(\$1,\$2,\$3,\$4\),\(\$5,\$6,\$7,\$8\)`).
				WithArgs(1, 1, 1, 0, 2, 1, 1, 100).
				WillReturnResult(sqlmock.NewResult(0, 2))

			// when
			err := repository.Create(transactions)

			// then
			assert.NoError(t, err)
		})
	})

	t.Run("FindByID", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			id := 1

			dbmock.
				ExpectQuery(`select id, chapterid, payinguserid, receivinguserid, amount from transactions where id = \$1`).
				WillReturnError(errors.New("database error"))

			// when
			transactions, err := repository.FindByID(id)

			// then
			assert.Error(t, err)
			assert.Nil(t, transactions)
		})

		t.Run("should return transactions by id", func(t *testing.T) {
			// given
			id := 1

			dbmock.
				ExpectQuery(`select id, chapterid, payinguserid, receivinguserid, amount from transactions where id = \$1`).
				WillReturnRows(sqlmock.NewRows([]string{"id", "chapterid", "payinguserid", "receivinguserid", "amount"}).
					AddRow(1, 1, 1, 1, 0))

			// when
			transactions, err := repository.FindByID(id)

			// then
			assert.NoError(t, err)
			assert.Len(t, transactions, 1)
		})
	})
}
