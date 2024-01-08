package products

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	_ "github.com/rqlite/gorqlite/stdlib"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products/model"
	"log"
)

type RQLiteRepository struct {
	db             *sql.DB
	productBuilder *sqlbuilder.Struct
}

const TableName = "product"

func NewRQLiteRepository(connectionString string) *RQLiteRepository {
	db, err := sql.Open("rqlite", connectionString)
	if err != nil {
		panic(fmt.Sprintf("Can't open repository: %v", err))
	}
	return &RQLiteRepository{
		db:             db,
		productBuilder: sqlbuilder.NewStruct(new(model.Product)).For(sqlbuilder.SQLite),
	}
}

func (r *RQLiteRepository) Create(product *model.Product) (*model.Product, error) {
	query, args := r.productBuilder.WithoutTag("pk").InsertInto("product", product).Build()

	transaction, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	result, err := transaction.Exec(query, args...)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		return nil, errors.New(ErrorProductAlreadyExists)
	}
	err = transaction.Commit()
	if err != nil {
		return nil, err
	}

	newId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	product.Id = uint64(newId)

	return product, nil
}

func (r *RQLiteRepository) FindAll() ([]*model.Product, error) {
	selectBuilder := r.productBuilder.SelectFrom(TableName)
	query, _ := selectBuilder.Build()

	transaction, err := r.db.Begin()
	if err != nil {
		return nil, errors.New(ErrorProductsList)
	}
	rows, err := transaction.Query(query)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		return nil, errors.New(ErrorProductsList)
	}

	var users = make([]*model.Product, 0)
	for rows.Next() {
		user := new(model.Product)
		err := rows.Scan(r.productBuilder.Addr(&user)...)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	err = transaction.Commit()
	if err != nil {
		return nil, errors.New(ErrorProductsList)
	}

	return users, nil
}

func (r *RQLiteRepository) FindById(id uint64) (*model.Product, error) {
	selectBuilder := r.productBuilder.SelectFrom(TableName)
	selectBuilder.Where(selectBuilder.Equal(TableName+".id", id))
	query, args := selectBuilder.Build()

	transaction, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	row := transaction.QueryRow(query, args...)

	user := &model.Product{}
	err = row.Scan(r.productBuilder.Addr(&user)...)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(ErrorProductNotFound)
		}
		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		return nil, errors.New(ErrorProductNotFound)
	}

	return user, nil
}

func (r *RQLiteRepository) FindByEan(ean string) (*model.Product, error) {
	selectBuilder := r.productBuilder.SelectFrom(TableName)
	selectBuilder.Where(selectBuilder.Equal(TableName+".ean", ean))
	query, args := selectBuilder.Build()

	transaction, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	row := transaction.QueryRow(query, args...)

	user := &model.Product{}
	err = row.Scan(r.productBuilder.Addr(&user)...)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(ErrorProductNotFound)
		}
		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		return nil, errors.New(ErrorProductNotFound)
	}

	return user, nil
}

func (r *RQLiteRepository) Update(product *model.Product) (*model.Product, error) {
	updateBuilder := r.productBuilder.Update(TableName, product)
	updateBuilder.Where(updateBuilder.Equal(TableName+".id", product.Id))
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
		return nil, errors.New(ErrorProductUpdate)
	}
	err = transaction.Commit()
	if err != nil {
		return nil, err
	}

	affectedRows, _ := result.RowsAffected()
	if affectedRows == 0 {
		return nil, errors.New(ErrorProductUpdate)
	}

	return product, nil
}

func (r *RQLiteRepository) Delete(product *model.Product) error {
	deleteBuilder := r.productBuilder.DeleteFrom(TableName)
	query, args := deleteBuilder.Where(deleteBuilder.Equal(TableName+".id", product.Id)).Build()

	transaction, err := r.db.Begin()
	if err != nil {
		return errors.New(ErrorProductDeletion)
	}
	result, err := transaction.Exec(query, args...)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		return errors.New(ErrorProductDeletion)
	}
	err = transaction.Commit()
	if err != nil {
		return err
	}

	affectedRows, _ := result.RowsAffected()
	if affectedRows == 0 {
		return errors.New(ErrorProductDeletion)
	}

	return nil
}
