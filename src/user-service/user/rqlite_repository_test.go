package user

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/huandu/go-sqlbuilder"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user/model"
	"reflect"
	"testing"
)

func TestRQLiteRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rqliteRepository := RQLiteRepository{
		db:          db,
		userBuilder: sqlbuilder.NewStruct(new(model.User)).For(sqlbuilder.SQLite),
	}

	user := model.User{
		Email:    "ada.lovelace@gmail.com",
		Password: []byte("123456"),
		Name:     "Ada Lovelace",
		Role:     model.Customer,
	}

	t.Run("Create user with success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO "+RQLiteTableName).
			WithArgs(user.Email, user.Password, user.Name, user.Role).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		_, err = rqliteRepository.Create(&user)
		if err != nil {
			t.Error(err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Can't create user with existing mail address", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO "+RQLiteTableName).
			WithArgs(user.Email, user.Password, user.Name, user.Role).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()

		_, err = rqliteRepository.Create(&user)
		if err.Error() != ErrorUserAlreadyExists {
			t.Error(err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Database error should return error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO "+RQLiteTableName).
			WithArgs(user.Email, user.Password, user.Name, user.Role).
			WillReturnError(errors.New("database has failed"))
		mock.ExpectRollback()

		_, err = rqliteRepository.Create(&user)
		if err == nil {
			t.Error("there should be an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})
}

func TestRQLiteRepository_FindAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rqliteRepository := RQLiteRepository{
		db:          db,
		userBuilder: sqlbuilder.NewStruct(new(model.User)).For(sqlbuilder.SQLite),
	}

	users := []*model.User{
		{
			Id:       1,
			Email:    "ada.lovelace@gmail.com",
			Password: []byte("123456"),
			Name:     "Ada Lovelace",
			Role:     model.Customer,
		},
		{
			Id:       2,
			Email:    "alan.turin@gmail.com",
			Password: []byte("123456"),
			Name:     "Alan Turing",
			Role:     model.Customer,
		},
	}

	t.Run("Successfully fetch all users", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT %[1]s.id, %[1]s.email, %[1]s.password, %[1]s.name, %[1]s.role FROM %[1]s`, RQLiteTableName)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "name", "role"}).
				AddRow(users[0].Id, users[0].Email, base64.StdEncoding.EncodeToString(users[0].Password), users[0].Name, users[0].Role).
				AddRow(users[1].Id, users[1].Email, base64.StdEncoding.EncodeToString(users[1].Password), users[1].Name, users[1].Role))
		mock.ExpectCommit()

		fetchedUsers, err := rqliteRepository.FindAll()

		if err != nil {
			t.Error("Can't fetch users")
		}

		if len(fetchedUsers) != len(users) {
			t.Errorf("Unexpected user count. Expected %d, got %d", len(users), len(fetchedUsers))
		}

		if !reflect.DeepEqual(users, fetchedUsers) {
			t.Error("Fetched users do not match expected users")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Database error should return error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT %[1]s.id, %[1]s.email, %[1]s.password, %[1]s.name, %[1]s.role FROM %[1]s`, RQLiteTableName)).
			WillReturnError(errors.New("database has failed"))
		mock.ExpectRollback()

		_, err := rqliteRepository.FindAll()
		if err == nil {
			t.Error("there should be an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})
}

func TestRQLiteRepository_FindAllByRole(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rqliteRepository := RQLiteRepository{
		db:          db,
		userBuilder: sqlbuilder.NewStruct(new(model.User)).For(sqlbuilder.SQLite),
	}

	users := []*model.User{
		{
			Id:       1,
			Email:    "ada.lovelace@gmail.com",
			Password: []byte("123456"),
			Name:     "Ada Lovelace",
			Role:     model.Merchant,
		},
		{
			Id:       2,
			Email:    "alan.turin@gmail.com",
			Password: []byte("123456"),
			Name:     "Alan Turing",
			Role:     model.Merchant,
		},
	}

	t.Run("Successfully fetch all merchant users", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT %[1]s.id, %[1]s.email, %[1]s.password, %[1]s.name, %[1]s.role FROM %[1]s WHERE %[1]s.role = \?`, RQLiteTableName)).
			WithArgs(model.Merchant).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "name", "role"}).
				AddRow(users[0].Id, users[0].Email, base64.StdEncoding.EncodeToString(users[0].Password), users[0].Name, users[0].Role).
				AddRow(users[1].Id, users[1].Email, base64.StdEncoding.EncodeToString(users[1].Password), users[1].Name, users[1].Role))
		mock.ExpectCommit()

		fetchedUsers, err := rqliteRepository.FindAllByRole(model.Merchant)

		if err != nil {
			t.Error("Can't fetch users")
		}

		if len(fetchedUsers) != len(users) {
			t.Errorf("Unexpected user count. Expected %d, got %d", len(users), len(fetchedUsers))
		}

		if !reflect.DeepEqual(users, fetchedUsers) {
			t.Error("Fetched users do not match expected users")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Database error should return error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT %[1]s.id, %[1]s.email, %[1]s.password, %[1]s.name, %[1]s.role FROM %[1]s WHERE %[1]s.role = \?`, RQLiteTableName)).
			WithArgs(model.Merchant).
			WillReturnError(errors.New("database has failed"))
		mock.ExpectRollback()

		_, err := rqliteRepository.FindAllByRole(model.Merchant)
		if err == nil {
			t.Error("there should be an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})
}

func TestRQLiteRepository_FindById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rqliteRepository := RQLiteRepository{
		db:          db,
		userBuilder: sqlbuilder.NewStruct(new(model.User)).For(sqlbuilder.SQLite),
	}

	user := model.User{
		Id:       1,
		Email:    "ada.lovelace@gmail.com",
		Password: []byte("123456"),
		Name:     "Ada Lovelace",
		Role:     model.Customer,
	}

	t.Run("Successfully fetch user", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT %[1]s.id, %[1]s.email, %[1]s.password, %[1]s.name, %[1]s.role FROM %[1]s WHERE %[1]s.id = \?`, RQLiteTableName)).
			WithArgs(user.Id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "name", "role"}).
				AddRow(user.Id, user.Email, base64.StdEncoding.EncodeToString(user.Password), user.Name, user.Role))
		mock.ExpectCommit()

		fetchedUser, err := rqliteRepository.FindById(user.Id)
		if err != nil {
			t.Errorf("Can't find expected user with id %d: %v", user.Id, err)
		}

		if !reflect.DeepEqual(user, *fetchedUser) {
			t.Error("Fetched user does not match original user")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Fail to fetch user (user not found)", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT %[1]s.id, %[1]s.email, %[1]s.password, %[1]s.name, %[1]s.role FROM %[1]s WHERE %[1]s.id = \?`, RQLiteTableName)).
			WithArgs(user.Id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "name", "role"}))
		mock.ExpectRollback()

		_, err := rqliteRepository.FindById(user.Id)
		if err == nil {
			t.Errorf("there should be an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Database error should return error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT %[1]s.id, %[1]s.email, %[1]s.password, %[1]s.name, %[1]s.role FROM %[1]s WHERE %[1]s.id = \?`, RQLiteTableName)).
			WithArgs(user.Id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "name", "role"}))
		mock.ExpectRollback()

		_, err := rqliteRepository.FindById(user.Id)
		if err == nil {
			t.Error("there should be an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})
}

func TestRQLiteRepository_FindByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rqliteRepository := RQLiteRepository{
		db:          db,
		userBuilder: sqlbuilder.NewStruct(new(model.User)).For(sqlbuilder.SQLite),
	}

	user := model.User{
		Id:       1,
		Email:    "ada.lovelace@gmail.com",
		Password: []byte("123456"),
		Name:     "Ada Lovelace",
		Role:     model.Customer,
	}

	t.Run("Successfully fetch user", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT %[1]s.id, %[1]s.email, %[1]s.password, %[1]s.name, %[1]s.role FROM %[1]s WHERE %[1]s.email = \?`, RQLiteTableName)).
			WithArgs(user.Email).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "name", "role"}).
				AddRow(user.Id, user.Email, base64.StdEncoding.EncodeToString(user.Password), user.Name, user.Role))
		mock.ExpectCommit()

		fetchedUser, err := rqliteRepository.FindByEmail(user.Email)
		if err != nil {
			t.Errorf("Can't find expected user with id %d: %v", user.Id, err)
		}

		if !reflect.DeepEqual(user, *fetchedUser) {
			t.Error("Fetched user does not match original user")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Fail to fetch user (user not found)", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT %[1]s.id, %[1]s.email, %[1]s.password, %[1]s.name, %[1]s.role FROM %[1]s WHERE %[1]s.email = \?`, RQLiteTableName)).
			WithArgs(user.Email).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "name", "role"}))
		mock.ExpectRollback()

		_, err := rqliteRepository.FindByEmail(user.Email)
		if err == nil {
			t.Errorf("there should be an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Database error should return error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT %[1]s.id, %[1]s.email, %[1]s.password, %[1]s.name, %[1]s.role FROM %[1]s WHERE %[1]s.email = \?`, RQLiteTableName)).
			WithArgs(user.Email).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "name", "role"}))
		mock.ExpectRollback()

		_, err := rqliteRepository.FindByEmail(user.Email)
		if err == nil {
			t.Error("there should be an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})
}

func TestRQLiteRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rqliteRepository := RQLiteRepository{
		db:          db,
		userBuilder: sqlbuilder.NewStruct(new(model.User)).For(sqlbuilder.SQLite),
	}

	changedUser := model.User{
		Id:       1,
		Email:    "ada.lovelace@gmail.com",
		Password: []byte("123456"),
		Name:     "Ada Lovelace",
		Role:     model.Customer,
	}

	t.Run("Update user with success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.
			ExpectExec(`UPDATE `+RQLiteTableName).
			WithArgs(changedUser.Id, changedUser.Email, changedUser.Password, changedUser.Name, changedUser.Role, changedUser.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		updatedUser, err := rqliteRepository.Update(&changedUser)
		if updatedUser.Name != updatedUser.Name || err != nil {
			t.Error("Failed to update user")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Update user with fail (user not found)", func(t *testing.T) {
		mock.ExpectBegin()
		mock.
			ExpectExec(`UPDATE `+RQLiteTableName).
			WithArgs(changedUser.Id, changedUser.Email, changedUser.Password, changedUser.Name, changedUser.Role, changedUser.Id).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		_, err := rqliteRepository.Update(&changedUser)
		if err == nil {
			t.Error("there should be an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Database error should return error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.
			ExpectExec(`UPDATE `+RQLiteTableName).
			WithArgs(changedUser.Id, changedUser.Email, changedUser.Password, changedUser.Name, changedUser.Role, changedUser.Id).
			WillReturnError(errors.New("database has failed"))
		mock.ExpectRollback()

		_, err = rqliteRepository.Update(&changedUser)
		if err == nil {
			t.Error("there should be an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

}

func TestRQLiteRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rqliteRepository := RQLiteRepository{
		db:          db,
		userBuilder: sqlbuilder.NewStruct(new(model.User)).For(sqlbuilder.SQLite),
	}

	userToDelete := model.User{
		Id:       1,
		Email:    "ada.lovelace@gmail.com",
		Password: []byte("123456"),
		Name:     "Ada Lovelace",
		Role:     model.Customer,
	}

	t.Run("Delete user with success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.
			ExpectExec(`DELETE FROM user WHERE user.id = \?`).
			WithArgs(userToDelete.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = rqliteRepository.Delete(&userToDelete)
		if err != nil {
			t.Errorf("Failed to delete user with id %d", userToDelete.Id)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Delete user with fail (user not found)", func(t *testing.T) {
		mock.ExpectBegin()
		mock.
			ExpectExec(fmt.Sprintf(`DELETE FROM %[1]s WHERE %[1]s.id = \?`, RQLiteTableName)).
			WithArgs(userToDelete.Id).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := rqliteRepository.Delete(&userToDelete)
		if err == nil {
			t.Error("there should be an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Database error should return error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.
			ExpectExec(fmt.Sprintf(`DELETE FROM %[1]s WHERE %[1]s.id = \?`, RQLiteTableName)).
			WithArgs(userToDelete.Id).
			WillReturnError(errors.New("database has failed"))
		mock.ExpectRollback()

		err = rqliteRepository.Delete(&userToDelete)
		if err == nil {
			t.Error("there should be an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})
}
