package user

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/user/model"
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
			users := []*model.DbUser{{
				ID:          1,
				Email:       "doesnt matter",
				Password:    []byte("doesnt matter"),
				Username:    "doesnt matter",
				ProfileName: "doesnt matter",
				Balance:     0,
			}}

			dbmock.
				ExpectExec(`insert into users`).
				WillReturnError(errors.New("database error"))

			// when
			err := repository.Create(users)

			// then
			assert.Error(t, err)
		})

		t.Run("should insert users in batches", func(t *testing.T) {
			// given
			users := []*model.DbUser{
				{
					ID:          1,
					Email:       "test@test.com",
					Password:    []byte("test"),
					Username:    "test",
					ProfileName: "Toni Tester",
					Balance:     0,
				},
				{
					ID:          2,
					Email:       "abc@abc.com",
					Password:    []byte("abc"),
					Username:    "abc",
					ProfileName: "ABC ABC",
					Balance:     0,
				},
			}

			dbmock.
				ExpectExec(`insert into users \(email, username, password, profile_name\) values \(\$1,\$2,\$3,\$4\),\(\$5,\$6,\$7,\$8\)`).
				WithArgs("test@test.com", "test", []byte("test"), "Toni Tester", "abc@abc.com", "abc", []byte("abc"), "ABC ABC").
				WillReturnResult(sqlmock.NewResult(0, 2))

			// when
			err := repository.Create(users)

			// then
			assert.NoError(t, err)
		})
	})

	t.Run("Update", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			user := &model.UpdateUser{
				ProfileName: "doesnt matter",
				Balance:     0,
			}

			dbmock.
				ExpectExec(`update users set profile_name = \$1, balance = \$2 where username = \$3`).
				WillReturnError(errors.New("database error"))

			// when
			err := repository.Update("tester", user)

			// then
			assert.Error(t, err)
		})

		t.Run("should update user", func(t *testing.T) {
			// given
			newUserData := &model.UpdateUser{
				ProfileName: "doesnt matter",
				Balance:     0,
			}

			dbmock.
				ExpectExec(`update users set profile_name = \$1, balance = \$2 where username = \$3`).
				WithArgs("doesnt matter", 0, "tester")

			// when
			err := repository.Update("tester", newUserData)

			// then
			assert.Error(t, err)
		})
	})

	t.Run("FindAll", func(t *testing.T) {
		t.Run("should return all products", func(t *testing.T) {
			// given
			dbmock.ExpectQuery(`select (.*) from users`).
				WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "username", "profile_name", "balance"}).
					AddRow(1, "test@test.com", []byte("hash"), "tester", "Toni Tester", 0).
					AddRow(2, "abc@abc.com", []byte("hash"), "abc", "ABC ABC", 0))

			// when
			products, err := repository.FindAll()

			// then
			assert.NoError(t, err)
			assert.NoError(t, dbmock.ExpectationsWereMet())
			assert.Len(t, products, 2)
			assert.Equal(t, "tester", products[0].Username)
			assert.Equal(t, "abc", products[1].Username)
		})
	})

	t.Run("FindByEmail", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			email := "test@test.com"

			dbmock.
				ExpectQuery(`select id, email, password, username, profile_name, balance from users where email = \$1`).
				WillReturnError(errors.New("database error"))

			// when
			users, err := repository.FindByEmail(email)

			// then
			assert.Error(t, err)
			assert.Nil(t, users)
		})

		t.Run("should return users by email", func(t *testing.T) {
			// given
			email := "test@test.com"

			dbmock.
				ExpectQuery(`select id, email, password, username, profile_name, balance from users where email = \$1`).
				WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "username", "profile_name", "balance"}).
					AddRow(1, "test@test.com", []byte("hash"), "tester", "Toni Tester", 0))

			// when
			users, err := repository.FindByEmail(email)

			// then
			assert.NoError(t, err)
			assert.Len(t, users, 1)
		})
	})

	t.Run("FindByUsername", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			username := "tester"

			dbmock.
				ExpectQuery(`select id, email, password, username, profile_name, balance from users where username = \$1`).
				WillReturnError(errors.New("database error"))

			// when
			users, err := repository.FindByUsername(username)

			// then
			assert.Error(t, err)
			assert.Nil(t, users)
		})

		t.Run("should return users by username", func(t *testing.T) {
			// given
			username := "tester"

			dbmock.
				ExpectQuery(`select id, email, password, username, profile_name, balance from users where username = \$1`).
				WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "username", "profile_name", "balance"}).
					AddRow(1, "test@test.com", []byte("hash"), "tester", "Toni Tester", 0))

			// when
			users, err := repository.FindByUsername(username)

			// then
			assert.NoError(t, err)
			assert.Len(t, users, 1)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			users := []*model.DbUser{
				{
					ID:          1,
					Email:       "test@test.com",
					Password:    []byte("hash"),
					Username:    "tester",
					ProfileName: "Toni Tester",
					Balance:     0,
				},
				{
					ID:          2,
					Email:       "abc@abc.com",
					Password:    []byte("hash"),
					Username:    "abc",
					ProfileName: "ABC ABC",
					Balance:     0,
				},
			}

			dbmock.
				ExpectExec(`delete from users`).
				WillReturnError(errors.New("database error"))

			// when
			err := repository.Delete(users)

			// then
			assert.Error(t, err)
		})

		t.Run("should delete users in batches", func(t *testing.T) {
			// given
			users := []*model.DbUser{
				{
					ID:          1,
					Email:       "test@test.com",
					Password:    []byte("hash"),
					Username:    "tester",
					ProfileName: "Toni Tester",
					Balance:     0,
				},
				{
					ID:          2,
					Email:       "abc@abc.com",
					Password:    []byte("hash"),
					Username:    "abc",
					ProfileName: "ABC ABC",
					Balance:     0,
				},
			}

			dbmock.
				ExpectExec(`delete from users where username in \(\$1,\$2\)`).
				WithArgs("tester", "abc").
				WillReturnResult(sqlmock.NewResult(0, 2))

			// when
			err := repository.Delete(users)

			// then
			assert.NoError(t, err)
		})
	})
}
