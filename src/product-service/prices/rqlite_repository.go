package prices

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	_ "github.com/rqlite/gorqlite/stdlib"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices/model"
	"log"
)

const RQLiteTableName = "price"

type RQLiteRepository struct {
	db           *sql.DB
	priceBuilder *sqlbuilder.Struct
}

func NewRQLiteRepository(connectionString string) *RQLiteRepository {
	db, err := sql.Open("rqlite", connectionString)
	if err != nil {
		panic(fmt.Sprintf("Can't open repository: %v", err))
	}
	return &RQLiteRepository{
		db:           db,
		priceBuilder: sqlbuilder.NewStruct(new(model.Price)).For(sqlbuilder.SQLite),
	}
}

func (r *RQLiteRepository) Create(price *model.Price) (*model.Price, error) {
	query, args := r.priceBuilder.InsertInto(RQLiteTableName, price).Build()

	transaction, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	_, err = transaction.Exec(query, args...)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		return nil, errors.New(ErrorPriceAlreadyExists)
	}
	err = transaction.Commit()
	if err != nil {
		return nil, err
	}

	return price, nil
}

func (r *RQLiteRepository) FindAll() ([]*model.Price, error) {
	selectBuilder := r.priceBuilder.SelectFrom(RQLiteTableName)
	query, _ := selectBuilder.Build()

	transaction, err := r.db.Begin()
	if err != nil {
		return nil, errors.New(ErrorPriceList)
	}
	rows, err := transaction.Query(query)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		return nil, errors.New(ErrorPriceList)
	}

	var prices = make([]*model.Price, 0)
	for rows.Next() {
		price := new(model.Price)
		err := rows.Scan(r.priceBuilder.Addr(&price)...)
		if err != nil {
			return nil, err
		}
		prices = append(prices, price)
	}

	err = transaction.Commit()
	if err != nil {
		return nil, errors.New(ErrorPriceList)
	}

	return prices, nil
}

func (r *RQLiteRepository) FindAllByUser(userId uint64) ([]*model.Price, error) {
	selectBuilder := r.priceBuilder.SelectFrom(RQLiteTableName)
	selectBuilder.Where(selectBuilder.Equal(RQLiteTableName+".userId", userId))
	query, args := selectBuilder.Build()

	transaction, err := r.db.Begin()
	if err != nil {
		return nil, errors.New(ErrorPriceList)
	}
	rows, err := transaction.Query(query, args...)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		return nil, errors.New(ErrorPriceList)
	}

	var prices = make([]*model.Price, 0)
	for rows.Next() {
		price := new(model.Price)
		err := rows.Scan(r.priceBuilder.Addr(&price)...)
		if err != nil {
			return nil, err
		}
		prices = append(prices, price)
	}

	err = transaction.Commit()
	if err != nil {
		return nil, errors.New(ErrorPriceList)
	}

	return prices, nil
}

func (r *RQLiteRepository) FindByIds(productId uint64, userId uint64) (*model.Price, error) {
	selectBuilder := r.priceBuilder.SelectFrom(RQLiteTableName)
	selectBuilder.Where(
		selectBuilder.Equal(RQLiteTableName+".userId", userId),
		selectBuilder.Equal(RQLiteTableName+".productId", productId))
	query, args := selectBuilder.Build()

	transaction, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	row := transaction.QueryRow(query, args...)

	price := &model.Price{}
	err = row.Scan(r.priceBuilder.Addr(&price)...)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(ErrorPriceNotFound)
		}
		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		return nil, errors.New(ErrorPriceNotFound)
	}

	return price, nil
}

func (r *RQLiteRepository) Update(price *model.Price) (*model.Price, error) {
	updateBuilder := r.priceBuilder.Update(RQLiteTableName, price)
	updateBuilder.Where(
		updateBuilder.Equal(RQLiteTableName+".userId", price.UserId),
		updateBuilder.Equal(RQLiteTableName+".productId", price.ProductId))
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
		return nil, errors.New(ErrorPriceUpdate)
	}
	err = transaction.Commit()
	if err != nil {
		return nil, err
	}

	affectedRows, _ := result.RowsAffected()
	if affectedRows == 0 {
		return nil, errors.New(ErrorPriceUpdate)
	}

	return price, nil
}

func (r *RQLiteRepository) Delete(price *model.Price) error {
	deleteBuilder := r.priceBuilder.DeleteFrom(RQLiteTableName)
	query, args := deleteBuilder.Where(
		deleteBuilder.Equal(RQLiteTableName+".userId", price.UserId),
		deleteBuilder.Equal(RQLiteTableName+".productId", price.ProductId)).Build()

	transaction, err := r.db.Begin()
	if err != nil {
		return errors.New(ErrorPriceDeletion)
	}
	result, err := transaction.Exec(query, args...)
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Println(err)
		}
		return errors.New(ErrorPriceDeletion)
	}
	err = transaction.Commit()
	if err != nil {
		return err
	}

	affectedRows, _ := result.RowsAffected()
	if affectedRows == 0 {
		return errors.New(ErrorPriceDeletion)
	}

	return nil
}
