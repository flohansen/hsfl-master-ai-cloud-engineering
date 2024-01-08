package userShoppingList

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/huandu/go-sqlbuilder"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList/model"
	"reflect"
	"testing"
)

func TestRQLiteRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rqliteRepository := RQLiteRepository{
		db:                  db,
		shoppinglistBuilder: sqlbuilder.NewStruct(new(model.UserShoppingList)).For(sqlbuilder.SQLite),
	}

	list := model.UserShoppingList{
		Id:          1,
		UserId:      1,
		Description: "New Shoppinglist",
		Completed:   false,
	}

	t.Run("Create list with success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO "+RQLiteTableName).
			WithArgs(list.UserId, list.Description, list.Completed).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		_, err = rqliteRepository.Create(&list)
		if err != nil {
			t.Error(err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Can't create lists with duplicate ean", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO "+RQLiteTableName).
			WithArgs(list.UserId, list.Description, list.Completed).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()

		_, err = rqliteRepository.Create(&list)
		if err.Error() != ErrorListAlreadyExists {
			t.Error(err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Database error should return error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO "+RQLiteTableName).
			WithArgs(list.UserId, list.Description, list.Completed).
			WillReturnError(errors.New("database has failed"))
		mock.ExpectRollback()

		_, err = rqliteRepository.Create(&list)
		if err == nil {
			t.Error("there should be an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})
}

func TestRQLiteRepository_FindAllById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rqliteRepository := RQLiteRepository{
		db:                  db,
		shoppinglistBuilder: sqlbuilder.NewStruct(new(model.UserShoppingList)).For(sqlbuilder.SQLite),
	}

	userId := uint64(2)

	lists := []*model.UserShoppingList{
		{
			Id:        1,
			UserId:    userId,
			Completed: false,
		},
		{
			Id:        2,
			UserId:    userId,
			Completed: false,
		},
		{
			Id:        3,
			UserId:    userId,
			Completed: true,
		},
	}

	t.Run("Successfully fetch all lists from user", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %[1]s WHERE %[1]s.userId = \?`, RQLiteTableName)).
			WithArgs(userId).
			WillReturnRows(sqlmock.NewRows([]string{"id", "userId", "description", "completed"}).
				AddRow(lists[0].Id, lists[0].UserId, lists[0].Description, lists[0].Completed).
				AddRow(lists[1].Id, lists[1].UserId, lists[1].Description, lists[1].Completed).
				AddRow(lists[2].Id, lists[2].UserId, lists[2].Description, lists[2].Completed))
		mock.ExpectCommit()

		fetchedLists, err := rqliteRepository.FindAllById(userId)

		if err != nil {
			t.Error(err)
		}

		if len(fetchedLists) != len(lists) {
			t.Errorf("Unexpected list count. Expected %d, got %d", len(lists), len(fetchedLists))
		}

		if !reflect.DeepEqual(lists, fetchedLists) {
			t.Error("Fetched lists do not match expected lists")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Database error should return error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %[1]s WHERE %[1]s.userId = \?`, RQLiteTableName)).
			WithArgs(userId).
			WillReturnError(errors.New("database has failed"))
		mock.ExpectRollback()

		_, err := rqliteRepository.FindAllById(userId)
		if err == nil {
			t.Error(err)
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
		db:                  db,
		shoppinglistBuilder: sqlbuilder.NewStruct(new(model.UserShoppingList)).For(sqlbuilder.SQLite),
	}

	list := model.UserShoppingList{
		Id:          1,
		UserId:      1,
		Description: "New List",
		Completed:   false,
	}

	t.Run("Successfully fetch list", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %[1]s WHERE %[1]s.id = \?`, RQLiteTableName)).
			WithArgs(list.Id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "userId", "description", "completed"}).
				AddRow(list.Id, list.UserId, list.Description, list.Completed))
		mock.ExpectCommit()

		fetchedList, err := rqliteRepository.FindById(list.Id)
		if err != nil {
			t.Errorf("Can't find expected list with id %d: %v", list.Id, err)
		}

		if !reflect.DeepEqual(list, *fetchedList) {
			t.Error("Fetched list does not match original list")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Fail to fetch list (list not found)", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %[1]s WHERE %[1]s.id = \?`, RQLiteTableName)).
			WithArgs(list.Id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "userId", "description", "completed"}))
		mock.ExpectRollback()

		_, err := rqliteRepository.FindById(list.Id)
		if err == nil {
			t.Errorf("there should be an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Database error should return error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %[1]s WHERE %[1]s.id = \?`, RQLiteTableName)).
			WithArgs(list.Id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "userId", "description", "completed"}))
		mock.ExpectRollback()

		_, err := rqliteRepository.FindById(list.Id)
		if err == nil {
			t.Error("there should be an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})
}

func TestRQLiteRepository_FindByIds(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rqliteRepository := RQLiteRepository{
		db:                  db,
		shoppinglistBuilder: sqlbuilder.NewStruct(new(model.UserShoppingList)).For(sqlbuilder.SQLite),
	}

	list := model.UserShoppingList{
		Id:          1,
		UserId:      1,
		Description: "New List",
		Completed:   false,
	}

	t.Run("Successfully fetch list", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %[1]s WHERE %[1]s.userId = \? AND %[1]s.id = \?`, RQLiteTableName)).
			WithArgs(list.UserId, list.Id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "userId", "description", "completed"}).
				AddRow(list.Id, list.UserId, list.Description, list.Completed))
		mock.ExpectCommit()

		fetchedList, err := rqliteRepository.FindByIds(list.Id, list.UserId)
		if err != nil {
			t.Errorf("Can't find expected list with id %d: %v", list.Id, err)
		}

		if !reflect.DeepEqual(list, *fetchedList) {
			t.Error("Fetched list does not match original list")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Fail to fetch list (list not found)", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %[1]s WHERE %[1]s.userId = \? AND %[1]s.id = \?`, RQLiteTableName)).
			WithArgs(list.UserId, list.Id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "userId", "description", "completed"}))
		mock.ExpectRollback()

		_, err := rqliteRepository.FindByIds(list.Id, list.UserId)
		if err == nil {
			t.Errorf("there should be an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Database error should return error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %[1]s WHERE %[1]s.userId = \? AND %[1]s.id = \?`, RQLiteTableName)).
			WithArgs(list.UserId, list.Id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "userId", "description", "completed"}))
		mock.ExpectRollback()

		_, err := rqliteRepository.FindByIds(list.Id, list.UserId)
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
		db:                  db,
		shoppinglistBuilder: sqlbuilder.NewStruct(new(model.UserShoppingList)).For(sqlbuilder.SQLite),
	}

	changedList := model.UserShoppingList{
		Id:          1,
		UserId:      1,
		Description: "Update list description",
		Completed:   true,
	}

	t.Run("Update list with success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE `+RQLiteTableName).
			WithArgs(changedList.Id, changedList.UserId, changedList.Description, changedList.Completed, changedList.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		updatedList, err := rqliteRepository.Update(&changedList)
		if reflect.DeepEqual(changedList, updatedList) || err != nil {
			t.Error("Failed to list list")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Update list with fail (list not found)", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE `+RQLiteTableName).
			WithArgs(changedList.Id, changedList.UserId, changedList.Description, changedList.Completed, changedList.Id).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		_, err := rqliteRepository.Update(&changedList)
		if err == nil {
			t.Error("there should be an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Database error should return error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE `+RQLiteTableName).
			WithArgs(changedList.Id, changedList.UserId, changedList.Description, changedList.Completed, changedList.Id).
			WillReturnError(errors.New("database has failed"))
		mock.ExpectRollback()

		_, err = rqliteRepository.Update(&changedList)
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
		db:                  db,
		shoppinglistBuilder: sqlbuilder.NewStruct(new(model.UserShoppingList)).For(sqlbuilder.SQLite),
	}

	listToDelete := model.UserShoppingList{
		Id:        1,
		UserId:    1,
		Completed: false,
	}

	t.Run("Delete list with success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.
			ExpectExec(fmt.Sprintf(`DELETE FROM %[1]s WHERE %[1]s.id = \?`, RQLiteTableName)).
			WithArgs(listToDelete.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = rqliteRepository.Delete(&listToDelete)
		if err != nil {
			t.Errorf("Failed to delete list with id %d", listToDelete.Id)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Delete list with fail (list not found)", func(t *testing.T) {
		mock.ExpectBegin()
		mock.
			ExpectExec(fmt.Sprintf(`DELETE FROM %[1]s WHERE %[1]s.id = \?`, RQLiteTableName)).
			WithArgs(listToDelete.Id).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := rqliteRepository.Delete(&listToDelete)
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
			WithArgs(listToDelete.Id).
			WillReturnError(errors.New("database has failed"))
		mock.ExpectRollback()

		err = rqliteRepository.Delete(&listToDelete)
		if err == nil {
			t.Error("there should be an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})
}
