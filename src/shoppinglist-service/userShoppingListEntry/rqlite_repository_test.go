package userShoppingListEntry

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/huandu/go-sqlbuilder"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingListEntry/model"
	"reflect"
	"testing"
)

func TestRQLiteRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rqliteRepository := RQLiteRepository{
		db:           db,
		entryBuilder: sqlbuilder.NewStruct(new(model.UserShoppingListEntry)).For(sqlbuilder.SQLite),
	}

	entry := model.UserShoppingListEntry{
		ShoppingListId: 1,
		ProductId:      1,
		Count:          2,
		Note:           "This is very important",
		Checked:        false,
	}

	t.Run("Create entry with success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO "+RQLiteTableName).
			WithArgs(entry.ShoppingListId, entry.ProductId, entry.Count, entry.Note, entry.Checked).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		_, err = rqliteRepository.Create(&entry)
		if err != nil {
			t.Error(err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Can't create prices with duplicate ean", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO "+RQLiteTableName).
			WithArgs(entry.ShoppingListId, entry.ProductId, entry.Count, entry.Note, entry.Checked).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()

		_, err = rqliteRepository.Create(&entry)
		if err.Error() != ErrorEntryAlreadyExists {
			t.Error(err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Database error should return error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO "+RQLiteTableName).
			WithArgs(entry.ShoppingListId, entry.ProductId, entry.Count, entry.Note, entry.Checked).
			WillReturnError(errors.New("database has failed"))
		mock.ExpectRollback()

		_, err = rqliteRepository.Create(&entry)
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
		db:           db,
		entryBuilder: sqlbuilder.NewStruct(new(model.UserShoppingListEntry)).For(sqlbuilder.SQLite),
	}

	shoppingListId := uint64(1)

	entries := []*model.UserShoppingListEntry{
		{
			ShoppingListId: shoppingListId,
			ProductId:      1,
			Count:          2,
			Note:           "This is very important",
			Checked:        false,
		},
		{
			ShoppingListId: shoppingListId,
			ProductId:      2,
			Count:          5,
			Note:           "I really want this",
			Checked:        false,
		},
		{
			ShoppingListId: shoppingListId,
			ProductId:      1,
			Count:          9,
			Note:           "Please get this",
			Checked:        false,
		},
	}

	t.Run("Successfully fetch all entries", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %[1]s WHERE %[1]s.shoppingListId = \?`, RQLiteTableName)).
			WithArgs(shoppingListId).
			WillReturnRows(sqlmock.NewRows([]string{"shoppingListId", "productId", "count", "note", "checked"}).
				AddRow(entries[0].ShoppingListId, entries[0].ProductId, entries[0].Count, entries[0].Note, entries[0].Checked).
				AddRow(entries[1].ShoppingListId, entries[1].ProductId, entries[1].Count, entries[1].Note, entries[1].Checked).
				AddRow(entries[2].ShoppingListId, entries[2].ProductId, entries[2].Count, entries[2].Note, entries[2].Checked))
		mock.ExpectCommit()

		fetchedEntries, err := rqliteRepository.FindAll(shoppingListId)

		if err != nil {
			t.Error(err)
		}

		if len(fetchedEntries) != len(entries) {
			t.Errorf("Unexpected entry count. Expected %d, got %d", len(entries), len(fetchedEntries))
		}

		if !reflect.DeepEqual(entries, fetchedEntries) {
			t.Error("Fetched entries do not match expected entries")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Database error should return error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %[1]s WHERE %[1]s.shoppingListId = \?`, RQLiteTableName)).
			WithArgs(shoppingListId).
			WillReturnError(errors.New("database has failed"))
		mock.ExpectRollback()

		_, err := rqliteRepository.FindAll(shoppingListId)
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
		db:           db,
		entryBuilder: sqlbuilder.NewStruct(new(model.UserShoppingListEntry)).For(sqlbuilder.SQLite),
	}

	entry := &model.UserShoppingListEntry{
		ShoppingListId: 1,
		ProductId:      1,
		Count:          2,
		Note:           "This is very important",
		Checked:        false,
	}

	t.Run("Successfully fetch entry", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %[1]s WHERE %[1]s.shoppingListId = \? AND %[1]s.productId = \?`, RQLiteTableName)).
			WithArgs(entry.ShoppingListId, entry.ProductId).
			WillReturnRows(sqlmock.NewRows([]string{"shoppingListId", "productId", "count", "note", "checked"}).
				AddRow(entry.ShoppingListId, entry.ProductId, entry.Count, entry.Note, entry.Checked))
		mock.ExpectCommit()

		fetchedEntry, err := rqliteRepository.FindByIds(entry.ShoppingListId, entry.ProductId)
		if err != nil {
			t.Errorf("Can't find expected entry with userId %d and productId %d: %v", entry.ShoppingListId, entry.ProductId, err)
		}

		if !reflect.DeepEqual(entry, fetchedEntry) {
			t.Error("Fetched entry does not match original entry")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Fail to fetch entry (entry not found)", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %[1]s WHERE %[1]s.shoppingListId = \? AND %[1]s.productId = \?`, RQLiteTableName)).
			WithArgs(entry.ShoppingListId, entry.ProductId).
			WillReturnRows(sqlmock.NewRows([]string{"shoppingListId", "productId", "count", "note", "checked"}))
		mock.ExpectRollback()

		_, err := rqliteRepository.FindByIds(entry.ShoppingListId, entry.ProductId)
		if err == nil {
			t.Errorf("there should be an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Database error should return error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %[1]s WHERE %[1]s.shoppingListId = \? AND %[1]s.productId = \?`, RQLiteTableName)).
			WithArgs(entry.ShoppingListId, entry.ProductId).
			WillReturnRows(sqlmock.NewRows([]string{"shoppingListId", "productId", "count", "note", "checked"}))
		mock.ExpectRollback()

		_, err := rqliteRepository.FindByIds(entry.ShoppingListId, entry.ProductId)
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
		db:           db,
		entryBuilder: sqlbuilder.NewStruct(new(model.UserShoppingListEntry)).For(sqlbuilder.SQLite),
	}

	changedEntry := model.UserShoppingListEntry{
		ShoppingListId: 1,
		ProductId:      1,
		Count:          3,
		Note:           "I really want this",
		Checked:        true,
	}

	t.Run("Update entry with success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(fmt.Sprintf(`UPDATE %[1]s SET .+ WHERE %[1]s.shoppingListId = \? AND %[1]s.productId = \?`, RQLiteTableName)).
			WithArgs(changedEntry.ShoppingListId,
				changedEntry.ProductId,
				changedEntry.Count,
				changedEntry.Note,
				changedEntry.Checked,
				changedEntry.ShoppingListId,
				changedEntry.ProductId).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		updatedEntry, err := rqliteRepository.Update(&changedEntry)
		if reflect.DeepEqual(changedEntry, updatedEntry) || err != nil {
			t.Error("Failed to update price")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Update entry with fail (entry not found)", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(fmt.Sprintf(`UPDATE %[1]s SET .+ WHERE %[1]s.shoppingListId = \? AND %[1]s.productId = \?`, RQLiteTableName)).
			WithArgs(changedEntry.ShoppingListId,
				changedEntry.ProductId,
				changedEntry.Count,
				changedEntry.Note,
				changedEntry.Checked,
				changedEntry.ShoppingListId,
				changedEntry.ProductId).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		_, err := rqliteRepository.Update(&changedEntry)
		if err == nil {
			t.Error("there should be an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Database error should return error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(fmt.Sprintf(`UPDATE %[1]s SET .+ WHERE %[1]s.shoppingListId = \? AND %[1]s.productId = \?`, RQLiteTableName)).
			WithArgs(changedEntry.ShoppingListId,
				changedEntry.ProductId,
				changedEntry.Count,
				changedEntry.Note,
				changedEntry.Checked,
				changedEntry.ShoppingListId,
				changedEntry.ProductId).
			WillReturnError(errors.New("database has failed"))
		mock.ExpectRollback()

		_, err = rqliteRepository.Update(&changedEntry)
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
		db:           db,
		entryBuilder: sqlbuilder.NewStruct(new(model.UserShoppingListEntry)).For(sqlbuilder.SQLite),
	}

	entryToDelete := model.UserShoppingListEntry{
		ShoppingListId: 1,
		ProductId:      1,
		Count:          1,
		Note:           "This is very important to me",
		Checked:        false,
	}

	t.Run("Delete price with success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.
			ExpectExec(fmt.Sprintf(`DELETE FROM %[1]s WHERE %[1]s.shoppingListId = \? AND %[1]s.productId = \?`, RQLiteTableName)).
			WithArgs(entryToDelete.ShoppingListId, entryToDelete.ProductId).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = rqliteRepository.Delete(&entryToDelete)
		if err != nil {
			t.Errorf("Failed to delete entry with shoppingListId %d and productId %d", entryToDelete.ShoppingListId, entryToDelete.ProductId)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Delete price with fail (price not found)", func(t *testing.T) {
		mock.ExpectBegin()
		mock.
			ExpectExec(fmt.Sprintf(`DELETE FROM %[1]s WHERE %[1]s.shoppingListId = \? AND %[1]s.productId = \?`, RQLiteTableName)).
			WithArgs(entryToDelete.ShoppingListId, entryToDelete.ProductId).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := rqliteRepository.Delete(&entryToDelete)
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
			ExpectExec(fmt.Sprintf(`DELETE FROM %[1]s WHERE %[1]s.shoppingListId = \? AND %[1]s.productId = \?`, RQLiteTableName)).
			WithArgs(entryToDelete.ShoppingListId, entryToDelete.ProductId).
			WillReturnError(errors.New("database has failed"))
		mock.ExpectRollback()

		err = rqliteRepository.Delete(&entryToDelete)
		if err == nil {
			t.Error("there should be an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})
}
