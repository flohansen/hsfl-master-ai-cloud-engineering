//go:build integration
// +build integration

package products

import (
	"context"
	_ "github.com/rqlite/gorqlite/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products/model"
	"reflect"
	"testing"
	"time"
)

const TestPort = "7003"

func TestIntegrationRQLiteRepository(t *testing.T) {
	container, err := prepareIntegrationTestRQLiteContainer()
	if err != nil {
		t.Error(err)
	}

	rqliteRepository := NewRQLiteRepository("http://localhost:" + TestPort + "/?disableClusterDiscovery=true")

	t.Run("TestIntegrationRQLiteRepository_Create", func(t *testing.T) {
		product := model.Product{
			Id:          1,
			Description: "Strauchtomaten",
			Ean:         "4014819040771",
		}

		t.Run("Create product with success", func(t *testing.T) {
			_, err = rqliteRepository.Create(&product)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("Can't create product with existing ean", func(t *testing.T) {
			_, err = rqliteRepository.Create(&product)
			if err.Error() != ErrorProductAlreadyExists {
				t.Error(err)
			}
		})

		t.Run("Database error should return error", func(t *testing.T) {
			_, err = rqliteRepository.Create(&product)
			if err == nil {
				t.Error("there should be an error")
			}
		})

		err := rqliteRepository.cleanTable()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("TestIntegrationRQLiteRepository_FindAll", func(t *testing.T) {
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

		for _, product := range products {
			_, err := rqliteRepository.Create(product)
			if err != nil {
				t.Error(err)
			}
		}

		t.Run("Successfully fetch all products", func(t *testing.T) {
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
		})

		err := rqliteRepository.cleanTable()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("TestIntegrationRQLiteRepository_FindByEan", func(t *testing.T) {
		product := model.Product{
			Id:          1,
			Description: "Strauchtomaten",
			Ean:         "4014819040771",
		}

		rqliteRepository.Create(&product)

		t.Run("Successfully fetch product", func(t *testing.T) {
			fetchedProduct, err := rqliteRepository.FindByEan(product.Ean)
			if err != nil {
				t.Errorf("Can't find expected product with id %d: %v", product.Id, err)
			}

			if !reflect.DeepEqual(product, *fetchedProduct) {
				t.Error("Fetched product does not match original product")
			}
		})

		t.Run("Fail to fetch product (product not found)", func(t *testing.T) {
			_, err := rqliteRepository.FindByEan("42")
			if err == nil {
				t.Errorf("there should be an error")
			}
		})

		err = rqliteRepository.cleanTable()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("TestIntegrationRQLiteRepository_FindById", func(t *testing.T) {
		product := model.Product{
			Id:          1,
			Description: "Strauchtomaten",
			Ean:         "4014819040771",
		}

		rqliteRepository.Create(&product)

		t.Run("Successfully fetch product", func(t *testing.T) {
			fetchedProduct, err := rqliteRepository.FindById(product.Id)
			if err != nil {
				t.Errorf("Can't find expected product with id %d: %v", product.Id, err)
			}

			if !reflect.DeepEqual(product, *fetchedProduct) {
				t.Error("Fetched product does not match original product")
			}
		})

		t.Run("Fail to fetch product (product not found)", func(t *testing.T) {
			_, err := rqliteRepository.FindById(2)
			if err == nil {
				t.Errorf("there should be an error")
			}
		})

		err = rqliteRepository.cleanTable()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("TestIntegrationRQLiteRepository_Update", func(t *testing.T) {
		product := model.Product{
			Id:          1,
			Description: "Strauchtomaten",
			Ean:         "4014819040771",
		}

		rqliteRepository.Create(&product)

		t.Run("Update product with success", func(t *testing.T) {
			changedProduct := model.Product{
				Id:          1,
				Description: "Sauftomaten",
				Ean:         "4014819040771",
			}

			updatedProduct, err := rqliteRepository.Update(&changedProduct)
			if updatedProduct.Description != changedProduct.Description || err != nil {
				t.Error("Failed to update product")
			}
		})

		t.Run("Update product with fail (product not found)", func(t *testing.T) {
			product := model.Product{
				Id:          55,
				Description: "Bier",
				Ean:         "12345484",
			}

			_, err := rqliteRepository.Update(&product)
			if err == nil {
				t.Error("there should be an error")
			}
		})

		err = rqliteRepository.cleanTable()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("TestIntegrationRQLiteRepository_Delete", func(t *testing.T) {
		productToDelete := model.Product{
			Id:          1,
			Description: "Strauchtomaten",
			Ean:         "4014819040771",
		}

		rqliteRepository.Create(&productToDelete)

		t.Run("Delete product with success", func(t *testing.T) {
			err = rqliteRepository.Delete(&productToDelete)
			if err != nil {
				t.Errorf("Failed to delete product with id %d", productToDelete.Id)
			}
		})

		t.Run("Delete product with fail (product not found)", func(t *testing.T) {
			nonExistingProduct := model.Product{
				Id:          1,
				Description: "Strauchtomaten",
				Ean:         "4014819040771",
			}

			err := rqliteRepository.Delete(&nonExistingProduct)
			if err == nil {
				t.Error("there should be an error")
			}
		})

		err = rqliteRepository.cleanTable()
		if err != nil {
			t.Error(err)
		}
	})

	t.Cleanup(func() {
		err = container.Stop(context.Background(), nil)
		if err != nil {
			return
		}
	})
}

func prepareIntegrationTestRQLiteContainer() (testcontainers.Container, error) {
	request := testcontainers.ContainerRequest{
		Image:        "rqlite/rqlite:8.15.0",
		ExposedPorts: []string{TestPort + ":4001/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("4001/tcp"),
			wait.ForLog(`.*HTTP API available at.*`).AsRegexp(),
			wait.ForLog(".*entering leader state.*").AsRegexp()).
			WithStartupTimeoutDefault(120 * time.Second).
			WithDeadline(360 * time.Second),
	}

	return testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: request,
			Started:          true,
		})
}
