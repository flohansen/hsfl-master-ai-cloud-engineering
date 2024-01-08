package prices

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/huandu/go-sqlbuilder"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices/model"
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
		priceBuilder: sqlbuilder.NewStruct(new(model.Price)).For(sqlbuilder.SQLite),
	}

	price := model.Price{
		UserId:    1,
		ProductId: 1,
		Price:     2.99,
	}

	t.Run("Create price with success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO "+RQLiteTableName).
			WithArgs(price.UserId, price.ProductId, price.Price).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		_, err = rqliteRepository.Create(&price)
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
			WithArgs(price.UserId, price.ProductId, price.Price).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()

		_, err = rqliteRepository.Create(&price)
		if err.Error() != ErrorPriceAlreadyExists {
			t.Error(err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Database error should return error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO "+RQLiteTableName).
			WithArgs(price.UserId, price.ProductId, price.Price).
			WillReturnError(errors.New("database has failed"))
		mock.ExpectRollback()

		_, err = rqliteRepository.Create(&price)
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
		priceBuilder: sqlbuilder.NewStruct(new(model.Price)).For(sqlbuilder.SQLite),
	}

	prices := []*model.Price{
		{
			UserId:    1,
			ProductId: 1,
			Price:     2.99,
		},
		{
			UserId:    2,
			ProductId: 3,
			Price:     0.55,
		},
	}

	t.Run("Successfully fetch all prices", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %s`, RQLiteTableName)).
			WillReturnRows(sqlmock.NewRows([]string{"userId", "productId", "price"}).
				AddRow(prices[0].UserId, prices[0].ProductId, prices[0].Price).
				AddRow(prices[1].UserId, prices[1].ProductId, prices[1].Price))
		mock.ExpectCommit()

		fetchedProducts, err := rqliteRepository.FindAll()

		if err != nil {
			t.Error("Can't fetch prices")
		}

		if len(fetchedProducts) != len(prices) {
			t.Errorf("Unexpected price count. Expected %d, got %d", len(prices), len(fetchedProducts))
		}

		if !reflect.DeepEqual(prices, fetchedProducts) {
			t.Error("Fetched prices do not match expected prices")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Database error should return error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %s`, RQLiteTableName)).
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

func TestRQLiteRepository_FindAllByUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rqliteRepository := RQLiteRepository{
		db:           db,
		priceBuilder: sqlbuilder.NewStruct(new(model.Price)).For(sqlbuilder.SQLite),
	}

	userIdMerchantA := uint64(1)

	pricesMerchantA := []*model.Price{
		{
			UserId:    userIdMerchantA,
			ProductId: 1,
			Price:     2.99,
		},
		{
			UserId:    userIdMerchantA,
			ProductId: 3,
			Price:     0.55,
		},
	}

	t.Run("Successfully fetch all prices from user", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %[1]s WHERE %[1]s.userId = ?`, RQLiteTableName)).
			WithArgs(userIdMerchantA).
			WillReturnRows(sqlmock.NewRows([]string{"userId", "productId", "price"}).
				AddRow(pricesMerchantA[0].UserId, pricesMerchantA[0].ProductId, pricesMerchantA[0].Price).
				AddRow(pricesMerchantA[1].UserId, pricesMerchantA[1].ProductId, pricesMerchantA[1].Price))
		mock.ExpectCommit()

		fetchedProducts, err := rqliteRepository.FindAllByUser(userIdMerchantA)

		if err != nil {
			t.Error("Can't fetch prices")
		}

		if len(fetchedProducts) != len(pricesMerchantA) {
			t.Errorf("Unexpected price count. Expected %d, got %d", len(pricesMerchantA), len(fetchedProducts))
		}

		if !reflect.DeepEqual(pricesMerchantA, fetchedProducts) {
			t.Error("Fetched prices do not match expected prices")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Database error should return error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %[1]s WHERE %[1]s.userId = ?`, RQLiteTableName)).
			WithArgs(userIdMerchantA).
			WillReturnError(errors.New("database has failed"))
		mock.ExpectRollback()

		_, err := rqliteRepository.FindAllByUser(userIdMerchantA)
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
		priceBuilder: sqlbuilder.NewStruct(new(model.Price)).For(sqlbuilder.SQLite),
	}

	price := model.Price{
		UserId:    1,
		ProductId: 1,
		Price:     2.99,
	}

	t.Run("Successfully fetch price", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %[1]s WHERE %[1]s.userId = \? AND %[1]s.productId = \?`, RQLiteTableName)).
			WithArgs(price.UserId, price.ProductId).
			WillReturnRows(sqlmock.NewRows([]string{"userId", "productId", "price"}).
				AddRow(price.UserId, price.ProductId, price.Price))
		mock.ExpectCommit()

		fetchedProduct, err := rqliteRepository.FindByIds(price.UserId, price.ProductId)
		if err != nil {
			t.Errorf("Can't find expected price with userId %d and productId %d: %v", price.UserId, price.ProductId, err)
		}

		if !reflect.DeepEqual(price, *fetchedProduct) {
			t.Error("Fetched price does not match original price")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Fail to fetch price (price not found)", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %[1]s WHERE %[1]s.userId = \? AND %[1]s.productId = \?`, RQLiteTableName)).
			WithArgs(price.UserId, price.ProductId).
			WillReturnRows(sqlmock.NewRows([]string{"userId", "productId", "price"}))
		mock.ExpectRollback()

		_, err := rqliteRepository.FindByIds(price.UserId, price.ProductId)
		if err == nil {
			t.Errorf("there should be an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Database error should return error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %[1]s WHERE %[1]s.userId = \? AND %[1]s.productId = \?`, RQLiteTableName)).
			WithArgs(price.UserId, price.ProductId).
			WillReturnRows(sqlmock.NewRows([]string{"userId", "productId", "price"}))
		mock.ExpectRollback()

		_, err := rqliteRepository.FindByIds(price.UserId, price.ProductId)
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
		priceBuilder: sqlbuilder.NewStruct(new(model.Price)).For(sqlbuilder.SQLite),
	}

	changedPrice := model.Price{
		UserId:    1,
		ProductId: 1,
		Price:     2.99,
	}

	t.Run("Update price with success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.
			ExpectExec(`UPDATE `+RQLiteTableName).
			WithArgs(changedPrice.UserId, changedPrice.ProductId, changedPrice.Price, changedPrice.UserId, changedPrice.ProductId).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		updatedProduct, err := rqliteRepository.Update(&changedPrice)
		if reflect.DeepEqual(changedPrice, updatedProduct) || err != nil {
			t.Error("Failed to update price")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Update price with fail (price not found)", func(t *testing.T) {
		mock.ExpectBegin()
		mock.
			ExpectExec(`UPDATE `+RQLiteTableName).
			WithArgs(changedPrice.UserId, changedPrice.ProductId, changedPrice.Price, changedPrice.UserId, changedPrice.ProductId).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		_, err := rqliteRepository.Update(&changedPrice)
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
			WithArgs(changedPrice.UserId, changedPrice.ProductId, changedPrice.Price, changedPrice.UserId, changedPrice.ProductId).
			WillReturnError(errors.New("database has failed"))
		mock.ExpectRollback()

		_, err = rqliteRepository.Update(&changedPrice)
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
		priceBuilder: sqlbuilder.NewStruct(new(model.Price)).For(sqlbuilder.SQLite),
	}

	priceToDelete := model.Price{
		UserId:    1,
		ProductId: 1,
		Price:     2.99,
	}

	t.Run("Delete price with success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.
			ExpectExec(fmt.Sprintf(`DELETE FROM %[1]s WHERE %[1]s.userId = \? AND %[1]s.productId = \?`, RQLiteTableName)).
			WithArgs(priceToDelete.UserId, priceToDelete.ProductId).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = rqliteRepository.Delete(&priceToDelete)
		if err != nil {
			t.Errorf("Failed to delete price with userId %d and productId %d", priceToDelete.UserId, priceToDelete.ProductId)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Delete price with fail (price not found)", func(t *testing.T) {
		mock.ExpectBegin()
		mock.
			ExpectExec(fmt.Sprintf(`DELETE FROM %[1]s WHERE %[1]s.userId = \? AND %[1]s.productId = \?`, RQLiteTableName)).
			WithArgs(priceToDelete.UserId, priceToDelete.ProductId).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := rqliteRepository.Delete(&priceToDelete)
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
			ExpectExec(fmt.Sprintf(`DELETE FROM %[1]s WHERE %[1]s.userId = \? AND %[1]s.productId = \?`, RQLiteTableName)).
			WithArgs(priceToDelete.UserId, priceToDelete.ProductId).
			WillReturnError(errors.New("database has failed"))
		mock.ExpectRollback()

		err = rqliteRepository.Delete(&priceToDelete)
		if err == nil {
			t.Error("there should be an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})
}
