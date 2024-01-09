package userShoppingListEntry

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	_ "github.com/rqlite/gorqlite/stdlib"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingListEntry/model"
	"log"
)

const (
	RQLiteTableName         = "shoppinglist_entry"
	RQLiteCreateTableQuery  = "CREATE TABLE IF NOT EXISTS " + RQLiteTableName + " ( shoppingListId BIGINT, productId BIGINT, count INTEGER NOT NULL, note TEXT, checked BOOLEAN NOT NULL, PRIMARY KEY (shoppingListId, productId) );"
	RQLiteCleanUpTableQuery = "DELETE FROM " + RQLiteTableName + ";"
)

type RQLiteRepository struct {
	db           *sql.DB
	entryBuilder *sqlbuilder.Struct
}

func NewRQLiteRepository(connectionString string) *RQLiteRepository {
	db, err := sql.Open("rqlite", connectionString)
	if err != nil {
		panic(fmt.Sprintf("Can't open repository: %v", err))
	}

	repository := &RQLiteRepository{
		db:           db,
		entryBuilder: sqlbuilder.NewStruct(new(model.UserShoppingListEntry)).For(sqlbuilder.SQLite),
	}

	err = repository.createTable()
	if err != nil {
		panic(err)
	}

	return repository
}

func (r *RQLiteRepository) Create(entry *model.UserShoppingListEntry) (*model.UserShoppingListEntry, error) {
	query, args := r.entryBuilder.InsertInto(RQLiteTableName, entry).Build()

	transaction, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	_, err = transaction.Exec(query, args...)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		return nil, errors.New(ErrorEntryAlreadyExists)
	}
	err = transaction.Commit()
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (r *RQLiteRepository) FindByIds(shoppingListId uint64, productId uint64) (*model.UserShoppingListEntry, error) {
	selectBuilder := r.entryBuilder.SelectFrom(RQLiteTableName)
	selectBuilder.Where(
		selectBuilder.Equal(RQLiteTableName+".shoppingListId", shoppingListId),
		selectBuilder.Equal(RQLiteTableName+".productId", productId))
	query, args := selectBuilder.Build()

	transaction, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	row := transaction.QueryRow(query, args...)

	entry := &model.UserShoppingListEntry{}
	err = row.Scan(r.entryBuilder.Addr(&entry)...)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(ErrorEntryNotFound)
		}
		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		return nil, errors.New(ErrorEntryNotFound)
	}

	return entry, nil
}

func (r *RQLiteRepository) FindAll(shoppingListId uint64) ([]*model.UserShoppingListEntry, error) {
	selectBuilder := r.entryBuilder.SelectFrom(RQLiteTableName)
	selectBuilder.Where(selectBuilder.Equal(RQLiteTableName+".shoppingListId", shoppingListId))
	query, args := selectBuilder.Build()

	transaction, err := r.db.Begin()
	if err != nil {
		return nil, errors.New(ErrorEntryNotFound)
	}
	rows, err := transaction.Query(query, args...)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		return nil, errors.New(ErrorEntryNotFound)
	}

	var entries = make([]*model.UserShoppingListEntry, 0)
	for rows.Next() {
		entry := new(model.UserShoppingListEntry)
		err := rows.Scan(r.entryBuilder.Addr(&entry)...)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	err = transaction.Commit()
	if err != nil {
		return nil, errors.New(ErrorEntryNotFound)
	}

	return entries, nil
}

func (r *RQLiteRepository) Update(entry *model.UserShoppingListEntry) (*model.UserShoppingListEntry, error) {
	updateBuilder := r.entryBuilder.Update(RQLiteTableName, entry)
	updateBuilder.Where(
		updateBuilder.Equal(RQLiteTableName+".shoppingListId", entry.ShoppingListId),
		updateBuilder.Equal(RQLiteTableName+".productId", entry.ProductId))
	query, args := updateBuilder.Build()

	transaction, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	result, err := transaction.Exec(query, args...)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		return nil, errors.New(ErrorEntryUpdate)
	}
	err = transaction.Commit()
	if err != nil {
		return nil, err
	}

	affectedRows, _ := result.RowsAffected()
	if affectedRows == 0 {
		return nil, errors.New(ErrorEntryUpdate)
	}

	return entry, nil
}

func (r *RQLiteRepository) Delete(entry *model.UserShoppingListEntry) error {
	deleteBuilder := r.entryBuilder.DeleteFrom(RQLiteTableName)
	query, args := deleteBuilder.Where(
		deleteBuilder.Equal(RQLiteTableName+".shoppingListId", entry.ShoppingListId),
		deleteBuilder.Equal(RQLiteTableName+".productId", entry.ProductId)).Build()

	transaction, err := r.db.Begin()
	if err != nil {
		return errors.New(ErrorEntryDeletion)
	}
	result, err := transaction.Exec(query, args...)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		return errors.New(ErrorEntryDeletion)
	}
	err = transaction.Commit()
	if err != nil {
		return err
	}

	affectedRows, _ := result.RowsAffected()
	if affectedRows == 0 {
		return errors.New(ErrorEntryDeletion)
	}

	return nil
}

func (r *RQLiteRepository) createTable() error {
	transaction, err := r.db.Begin()
	_, err = transaction.Exec(RQLiteCreateTableQuery)
	if err != nil {
		return err
	}
	err = transaction.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *RQLiteRepository) cleanTable() error {
	transaction, err := r.db.Begin()
	_, err = transaction.Exec(RQLiteCleanUpTableQuery)
	if err != nil {
		return err
	}
	err = transaction.Commit()
	if err != nil {
		return err
	}
	return nil
}
