package products

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/huandu/go-sqlbuilder"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products/model"
	"reflect"
	"testing"
)

func TestRQLiteRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rqliteRepository := RQLiteRepository{
		db:             db,
		productBuilder: sqlbuilder.NewStruct(new(model.Product)).For(sqlbuilder.SQLite),
	}

	product := model.Product{
		Id:          1,
		Description: "Strauchtomaten",
		Ean:         "4014819040771",
	}

	t.Run("Create product with success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO "+RQLiteTableName).
			WithArgs(product.Description, product.Ean).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		_, err = rqliteRepository.Create(&product)
		if err != nil {
			t.Error(err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Can't create products with duplicate ean", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO "+RQLiteTableName).
			WithArgs(product.Description, product.Ean).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()

		_, err = rqliteRepository.Create(&product)
		if err.Error() != ErrorProductAlreadyExists {
			t.Error(err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Database error should return error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO "+RQLiteTableName).
			WithArgs(product.Description, product.Ean).
			WillReturnError(errors.New("database has failed"))
		mock.ExpectRollback()

		_, err = rqliteRepository.Create(&product)
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
		db:             db,
		productBuilder: sqlbuilder.NewStruct(new(model.Product)).For(sqlbuilder.SQLite),
	}

	products := []*model.Product{
		{
			Id:          1,
			Description: "Strauchtomaten",
			Ean:         "4014819040771",
		},
		{
			Id:          2,
			Description: "Lauchzwiebeln",
			Ean:         "5001819040871",
		},
	}

	t.Run("Successfully fetch all products", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %s`, RQLiteTableName)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "description", "ean"}).
				AddRow(products[0].Id, products[0].Description, products[0].Ean).
				AddRow(products[1].Id, products[1].Description, products[1].Ean))
		mock.ExpectCommit()

		fetchedProducts, err := rqliteRepository.FindAll()

		if err != nil {
			t.Error("Can't fetch products")
		}

		if len(fetchedProducts) != len(products) {
			t.Errorf("Unexpected product count. Expected %d, got %d", len(products), len(fetchedProducts))
		}

		if !reflect.DeepEqual(products, fetchedProducts) {
			t.Error("Fetched products do not match expected products")
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

func TestRQLiteRepository_FindByEan(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rqliteRepository := RQLiteRepository{
		db:             db,
		productBuilder: sqlbuilder.NewStruct(new(model.Product)).For(sqlbuilder.SQLite),
	}

	product := model.Product{
		Id:          1,
		Description: "Strauchtomaten",
		Ean:         "4014819040771",
	}

	t.Run("Successfully fetch product", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %[1]s WHERE %[1]s.ean = \?`, RQLiteTableName)).
			WithArgs(product.Ean).
			WillReturnRows(sqlmock.NewRows([]string{"id", "description", "password"}).
				AddRow(product.Id, product.Description, product.Ean))
		mock.ExpectCommit()

		fetchedProduct, err := rqliteRepository.FindByEan(product.Ean)
		if err != nil {
			t.Errorf("Can't find expected product with ean %s: %v", product.Ean, err)
		}

		if !reflect.DeepEqual(product, *fetchedProduct) {
			t.Error("Fetched product does not match original product")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Fail to fetch product (product not found)", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %[1]s WHERE %[1]s.ean = \?`, RQLiteTableName)).
			WithArgs(product.Ean).
			WillReturnRows(sqlmock.NewRows([]string{"id", "description", "password"}))
		mock.ExpectRollback()

		_, err := rqliteRepository.FindByEan(product.Ean)
		if err == nil {
			t.Errorf("there should be an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Database error should return error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %[1]s WHERE %[1]s.ean = \?`, RQLiteTableName)).
			WithArgs(product.Ean).
			WillReturnRows(sqlmock.NewRows([]string{"id", "description", "password"}))
		mock.ExpectRollback()

		_, err := rqliteRepository.FindByEan(product.Ean)
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
		db:             db,
		productBuilder: sqlbuilder.NewStruct(new(model.Product)).For(sqlbuilder.SQLite),
	}

	product := model.Product{
		Id:          1,
		Description: "Strauchtomaten",
		Ean:         "4014819040771",
	}

	t.Run("Successfully fetch product", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %[1]s WHERE %[1]s.id = \?`, RQLiteTableName)).
			WithArgs(product.Id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "description", "ean"}).
				AddRow(product.Id, product.Description, product.Ean))
		mock.ExpectCommit()

		fetchedProduct, err := rqliteRepository.FindById(product.Id)
		if err != nil {
			t.Errorf("Can't find expected product with id %d: %v", product.Id, err)
		}

		if !reflect.DeepEqual(product, *fetchedProduct) {
			t.Error("Fetched product does not match original product")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Fail to fetch product (product not found)", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(fmt.Sprintf(`SELECT (.+) FROM %[1]s WHERE %[1]s.id = \?`, RQLiteTableName)).
			WithArgs(product.Id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "description", "ean"}))
		mock.ExpectRollback()

		_, err := rqliteRepository.FindById(product.Id)
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
			WithArgs(product.Id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "description", "ean"}))
		mock.ExpectRollback()

		_, err := rqliteRepository.FindById(product.Id)
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
		db:             db,
		productBuilder: sqlbuilder.NewStruct(new(model.Product)).For(sqlbuilder.SQLite),
	}

	changedProduct := model.Product{
		Id:          1,
		Description: "Strauchtomaten",
		Ean:         "4014819040771",
	}

	t.Run("Update product with success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.
			ExpectExec(`UPDATE `+RQLiteTableName).
			WithArgs(changedProduct.Id, changedProduct.Description, changedProduct.Ean, changedProduct.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		updatedProduct, err := rqliteRepository.Update(&changedProduct)
		if reflect.DeepEqual(changedProduct, updatedProduct) || err != nil {
			t.Error("Failed to product product")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Update product with fail (product not found)", func(t *testing.T) {
		mock.ExpectBegin()
		mock.
			ExpectExec(`UPDATE `+RQLiteTableName).
			WithArgs(changedProduct.Id, changedProduct.Description, changedProduct.Ean, changedProduct.Id).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		_, err := rqliteRepository.Update(&changedProduct)
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
			WithArgs(changedProduct.Id, changedProduct.Description, changedProduct.Ean, changedProduct.Id).
			WillReturnError(errors.New("database has failed"))
		mock.ExpectRollback()

		_, err = rqliteRepository.Update(&changedProduct)
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
		db:             db,
		productBuilder: sqlbuilder.NewStruct(new(model.Product)).For(sqlbuilder.SQLite),
	}

	productToDelete := model.Product{
		Id:          1,
		Description: "Strauchtomaten",
		Ean:         "4014819040771",
	}

	t.Run("Delete product with success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.
			ExpectExec(fmt.Sprintf(`DELETE FROM %[1]s WHERE %[1]s.id = \?`, RQLiteTableName)).
			WithArgs(productToDelete.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = rqliteRepository.Delete(&productToDelete)
		if err != nil {
			t.Errorf("Failed to delete product with id %d", productToDelete.Id)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Delete product with fail (product not found)", func(t *testing.T) {
		mock.ExpectBegin()
		mock.
			ExpectExec(fmt.Sprintf(`DELETE FROM %[1]s WHERE %[1]s.id = \?`, RQLiteTableName)).
			WithArgs(productToDelete.Id).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := rqliteRepository.Delete(&productToDelete)
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
			WithArgs(productToDelete.Id).
			WillReturnError(errors.New("database has failed"))
		mock.ExpectRollback()

		err = rqliteRepository.Delete(&productToDelete)
		if err == nil {
			t.Error("there should be an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})
}
