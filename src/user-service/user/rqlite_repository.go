package user

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	_ "github.com/rqlite/gorqlite/stdlib"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user/model"
	"log"
)

const (
	RQLiteTableName         = "user"
	RQLiteCreateTableQuery  = "CREATE TABLE IF NOT EXISTS " + RQLiteTableName + " ( id INTEGER PRIMARY KEY, email VARCHAR(255) UNIQUE, password BLOB, name VARCHAR(255), role BIGINT );"
	RQLiteCleanUpTableQuery = "DELETE FROM " + RQLiteTableName + ";"
)

type RQLiteRepository struct {
	db          *sql.DB
	userBuilder *sqlbuilder.Struct
}

func NewRQLiteRepository(connectionString string) *RQLiteRepository {
	db, err := sql.Open("rqlite", connectionString)
	if err != nil {
		panic(fmt.Sprintf("Can't open repository: %v", err))
	}

	repository := &RQLiteRepository{
		db:          db,
		userBuilder: sqlbuilder.NewStruct(new(model.User)).For(sqlbuilder.SQLite),
	}

	err = repository.createTable()
	if err != nil {
		panic(err)
	}

	return repository
}

func (r *RQLiteRepository) Create(user *model.User) (*model.User, error) {
	query, args := r.userBuilder.WithoutTag("pk").InsertInto(RQLiteTableName, user).Build()

	transaction, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	result, err := transaction.Exec(query, args...)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		return nil, errors.New(ErrorUserAlreadyExists)
	}
	err = transaction.Commit()
	if err != nil {
		return nil, err
	}

	newId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	user.Id = uint64(newId)

	return user, nil
}

func (r *RQLiteRepository) FindAll() ([]*model.User, error) {
	selectBuilder := r.userBuilder.SelectFrom(RQLiteTableName)
	query, _ := selectBuilder.Build()

	transaction, err := r.db.Begin()
	if err != nil {
		return nil, errors.New(ErrorUserList)
	}
	rows, err := transaction.Query(query)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		return nil, errors.New(ErrorUserList)
	}

	var users = make([]*model.User, 0)
	for rows.Next() {
		user := new(model.User)
		err := rows.Scan(r.userBuilder.Addr(&user)...)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	err = transaction.Commit()
	if err != nil {
		return nil, errors.New(ErrorUserList)
	}

	decodePassword(users...)

	return users, nil
}

func (r *RQLiteRepository) FindAllByRole(role model.Role) ([]*model.User, error) {
	selectBuilder := r.userBuilder.SelectFrom(RQLiteTableName)
	selectBuilder.Where(selectBuilder.Equal(RQLiteTableName+".role", role))
	query, args := selectBuilder.Build()

	transaction, err := r.db.Begin()
	rows, err := transaction.Query(query, args...)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		return nil, errors.New(ErrorUserList)
	}

	var users = make([]*model.User, 0)
	for rows.Next() {
		user := new(model.User)
		err := rows.Scan(r.userBuilder.Addr(&user)...)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	err = transaction.Commit()
	if err != nil {
		return nil, errors.New(ErrorUserList)
	}

	decodePassword(users...)

	return users, nil
}

func (r *RQLiteRepository) FindByEmail(email string) (*model.User, error) {
	selectBuilder := r.userBuilder.SelectFrom(RQLiteTableName)
	selectBuilder.Where(selectBuilder.Equal(RQLiteTableName+".email", email))
	query, args := selectBuilder.Build()

	transaction, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	row := transaction.QueryRow(query, args...)

	user := &model.User{}
	err = row.Scan(r.userBuilder.Addr(&user)...)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(ErrorUserNotFound)
		}
		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		return nil, errors.New(ErrorUserList)
	}

	decodePassword(user)

	return user, nil
}

func (r *RQLiteRepository) FindById(id uint64) (*model.User, error) {
	selectBuilder := r.userBuilder.SelectFrom(RQLiteTableName)
	selectBuilder.Where(selectBuilder.Equal(RQLiteTableName+".id", id))
	query, args := selectBuilder.Build()

	transaction, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	row := transaction.QueryRow(query, args...)

	user := &model.User{}
	err = row.Scan(r.userBuilder.Addr(&user)...)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(ErrorUserNotFound)
		}
		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		return nil, errors.New(ErrorUserNotFound)
	}

	decodePassword(user)

	return user, nil
}

func (r *RQLiteRepository) Update(user *model.User) (*model.User, error) {
	updateBuilder := r.userBuilder.Update(RQLiteTableName, user)
	updateBuilder.Where(updateBuilder.Equal(RQLiteTableName+".id", user.Id))
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
		return nil, errors.New(ErrorUserUpdate)
	}
	err = transaction.Commit()
	if err != nil {
		return nil, err
	}

	affectedRows, _ := result.RowsAffected()
	if affectedRows == 0 {
		return nil, errors.New(ErrorUserUpdate)
	}

	return user, nil
}

func (r *RQLiteRepository) Delete(user *model.User) error {
	deleteBuilder := r.userBuilder.DeleteFrom(RQLiteTableName)
	query, args := deleteBuilder.Where(deleteBuilder.Equal(RQLiteTableName+".id", user.Id)).Build()

	transaction, err := r.db.Begin()
	if err != nil {
		return errors.New(ErrorUserDeletion)
	}
	result, err := transaction.Exec(query, args...)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		return errors.New(ErrorUserDeletion)
	}
	err = transaction.Commit()
	if err != nil {
		return err
	}

	affectedRows, _ := result.RowsAffected()
	if affectedRows == 0 {
		return errors.New(ErrorUserUpdate)
	}

	return nil
}

func decodePassword(users ...*model.User) {
	for _, user := range users {
		decodedPassword, _ := base64.StdEncoding.DecodeString(string(user.Password))
		user.Password = decodedPassword
	}
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
