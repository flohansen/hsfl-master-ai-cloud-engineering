package userShoppingList

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/shoppinglist-service/userShoppingList/model"
	"log"
)

const (
	RQLiteTableName         = "shoppinglist"
	RQLiteCreateTableQuery  = "CREATE TABLE IF NOT EXISTS " + RQLiteTableName + " ( id INTEGER PRIMARY KEY, userId BIGINT NOT NULL, description VARCHAR(255), completed BOOLEAN NOT NULL );"
	RQLiteCleanUpTableQuery = "DELETE FROM " + RQLiteTableName + ";"
)

type RQLiteRepository struct {
	db                  *sql.DB
	shoppinglistBuilder *sqlbuilder.Struct
}

func NewRQLiteRepository(connectionString string) *RQLiteRepository {
	db, err := sql.Open("rqlite", connectionString)
	if err != nil {
		panic(fmt.Sprintf("Can't open repository: %v", err))
	}

	repository := &RQLiteRepository{
		db:                  db,
		shoppinglistBuilder: sqlbuilder.NewStruct(new(model.UserShoppingList)).For(sqlbuilder.SQLite),
	}

	err = repository.createTable()
	if err != nil {
		panic(err)
	}

	return repository
}

func (r *RQLiteRepository) Create(list *model.UserShoppingList) (*model.UserShoppingList, error) {
	query, args := r.shoppinglistBuilder.WithoutTag("pk").InsertInto(RQLiteTableName, list).Build()

	transaction, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	result, err := transaction.Exec(query, args...)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		return nil, errors.New(ErrorListAlreadyExists)
	}
	err = transaction.Commit()
	if err != nil {
		return nil, err
	}

	newId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	list.Id = uint64(newId)

	return list, nil
}

func (r *RQLiteRepository) FindAllById(userId uint64) ([]*model.UserShoppingList, error) {
	selectBuilder := r.shoppinglistBuilder.SelectFrom(RQLiteTableName)
	selectBuilder.Where(selectBuilder.Equal(RQLiteTableName+".userId", userId))
	query, args := selectBuilder.Build()

	transaction, err := r.db.Begin()
	if err != nil {
		return nil, errors.New(ErrorListNotFound)
	}
	rows, err := transaction.Query(query, args...)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		return nil, errors.New(ErrorListNotFound)
	}

	var lists = make([]*model.UserShoppingList, 0)
	for rows.Next() {
		list := new(model.UserShoppingList)
		err := rows.Scan(r.shoppinglistBuilder.Addr(&list)...)
		if err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}

	err = transaction.Commit()
	if err != nil {
		return nil, errors.New(ErrorListNotFound)
	}

	return lists, nil
}

func (r *RQLiteRepository) FindById(listId uint64) (*model.UserShoppingList, error) {
	selectBuilder := r.shoppinglistBuilder.SelectFrom(RQLiteTableName)
	selectBuilder.Where(selectBuilder.Equal(RQLiteTableName+".id", listId))
	query, args := selectBuilder.Build()

	transaction, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	row := transaction.QueryRow(query, args...)

	list := &model.UserShoppingList{}
	err = row.Scan(r.shoppinglistBuilder.Addr(&list)...)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(ErrorListNotFound)
		}
		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		return nil, errors.New(ErrorListNotFound)
	}

	return list, nil
}

func (r *RQLiteRepository) FindByIds(userId uint64, listId uint64) (*model.UserShoppingList, error) {
	selectBuilder := r.shoppinglistBuilder.SelectFrom(RQLiteTableName)
	selectBuilder.Where(
		selectBuilder.Equal(RQLiteTableName+".userId", userId),
		selectBuilder.Equal(RQLiteTableName+".id", listId))
	query, args := selectBuilder.Build()

	transaction, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	row := transaction.QueryRow(query, args...)

	list := &model.UserShoppingList{}
	err = row.Scan(r.shoppinglistBuilder.Addr(&list)...)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(ErrorListNotFound)
		}
		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		return nil, errors.New(ErrorListNotFound)
	}

	return list, nil
}

func (r *RQLiteRepository) Update(list *model.UserShoppingList) (*model.UserShoppingList, error) {
	updateBuilder := r.shoppinglistBuilder.Update(RQLiteTableName, list)
	updateBuilder.Where(updateBuilder.Equal(RQLiteTableName+".id", list.Id))
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
		return nil, errors.New(ErrorListUpdate)
	}
	err = transaction.Commit()
	if err != nil {
		return nil, err
	}

	affectedRows, _ := result.RowsAffected()
	if affectedRows == 0 {
		return nil, errors.New(ErrorListUpdate)
	}

	return list, nil
}

func (r *RQLiteRepository) Delete(list *model.UserShoppingList) error {
	deleteBuilder := r.shoppinglistBuilder.DeleteFrom(RQLiteTableName)
	query, args := deleteBuilder.Where(deleteBuilder.Equal(RQLiteTableName+".id", list.Id)).Build()

	transaction, err := r.db.Begin()
	if err != nil {
		return errors.New(ErrorListDeletion)
	}
	result, err := transaction.Exec(query, args...)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		return errors.New(ErrorListDeletion)
	}
	err = transaction.Commit()
	if err != nil {
		return err
	}

	affectedRows, _ := result.RowsAffected()
	if affectedRows == 0 {
		return errors.New(ErrorListDeletion)
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
