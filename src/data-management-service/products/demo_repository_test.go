package products

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/data-management-service/products/model"
	"reflect"
	"testing"
)

func TestNewDemoRepository(t *testing.T) {
	t.Run("Demo repository correct initialized", func(t *testing.T) {
		want := make(map[uint64]*model.Product)
		if got := NewDemoRepository(); !reflect.DeepEqual(got.products, want) {
			t.Errorf("NewDemoRepository().products = %v, want %v", got.products, want)
		}
	})
}

func TestDemoRepository_Create(t *testing.T) {
	// Prepare test
	demoRepository := NewDemoRepository()

	product := model.Product{
		Id:          1,
		Description: "Strauchtomaten",
		Ean:         4014819040771,
	}

	t.Run("Create product with success", func(t *testing.T) {
		_, err := demoRepository.Create(&product)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Check for doublet", func(t *testing.T) {
		_, err := demoRepository.Create(&product)
		if err.Error() != ErrorProductAlreadyExists {
			t.Error(err)
		}
	})
}

func TestDemoRepository_FindAll(t *testing.T) {
	// Prepare test
	demoRepository := NewDemoRepository()

	products := []*model.Product{
		{
			Id:          1,
			Description: "Strauchtomaten",
			Ean:         4014819040771,
		},
		{
			Id:          2,
			Description: "Lauchzwiebeln",
			Ean:         5001819040871,
		},
	}

	for _, product := range products {
		_, err := demoRepository.Create(product)
		if err != nil {
			t.Error("Failed to add prepared product for test")
		}
	}

	t.Run("Fetch all products", func(t *testing.T) {
		fetchedProducts, err := demoRepository.FindAll()
		if err != nil {
			t.Error("Can't fetch products")
		}

		if len(fetchedProducts) != len(products) {
			t.Errorf("Unexpected product count. Expected %d, got %d", len(products), len(fetchedProducts))
		}
	})

	productTests := []struct {
		name string
		want *model.Product
	}{
		{
			name: "first",
			want: products[0],
		},
		{
			name: "second",
			want: products[1],
		},
	}

	for i, tt := range productTests {
		t.Run("Is fetched product matching with "+tt.name+" added product?", func(t *testing.T) {
			fetchedProducts, _ := demoRepository.FindAll()
			if !reflect.DeepEqual(tt.want, fetchedProducts[i]) {
				t.Error("Fetched product does not match original product")
			}
		})
	}
}

func TestDemoRepository_FindById(t *testing.T) {
	// Prepare test
	demoRepository := NewDemoRepository()

	product := model.Product{
		Id:          1,
		Description: "Strauchtomaten",
		Ean:         4014819040771,
	}

	_, err := demoRepository.Create(&product)
	if err != nil {
		t.Fatal("Failed to add prepare product for test")
	}

	t.Run("Fetch product with existing id", func(t *testing.T) {
		_, err := demoRepository.FindById(product.Id)
		if err != nil {
			t.Errorf("Can't find expected product with id %d", product.Id)
		}

		t.Run("Is fetched product matching with added product?", func(t *testing.T) {
			fetchedProduct, _ := demoRepository.FindById(product.Id)
			if !reflect.DeepEqual(product, *fetchedProduct) {
				t.Error("Fetched product does not match original product")
			}
		})
	})

	t.Run("Non-existing product test", func(t *testing.T) {
		_, err = demoRepository.FindById(42)
		if err.Error() != ErrorProductNotFound {
			t.Error(err)
		}
	})
}

func TestDemoRepository_Update(t *testing.T) {
	// Prepare test
	demoRepository := NewDemoRepository()

	product := model.Product{
		Id:          1,
		Description: "Strauchtomaten",
		Ean:         4014819040771,
	}

	fetchedProduct, err := demoRepository.Create(&product)
	if err != nil {
		t.Error("Failed to add prepare product for test")
	}

	t.Run("Check if updated product has updated description", func(t *testing.T) {
		updateProduct := model.Product{
			Id:          1,
			Description: "Wittenseer Mineralwasser",
			Ean:         4014819040771,
		}
		updatedProduct, err := demoRepository.Update(&updateProduct)
		if err != nil {
			t.Error(err.Error())
		}

		if fetchedProduct.Description != updatedProduct.Description {
			t.Errorf("Failed to update product description. Got %s, want %s.",
				fetchedProduct.Description, updateProduct.Description)
		}
	})

}

func TestDemoRepository_Delete(t *testing.T) {
	// Prepare test
	productsRepository := NewDemoRepository()

	product := model.Product{
		Id:          1,
		Description: "Strauchtomaten",
		Ean:         4014819040771,
	}

	fetchedProduct, err := productsRepository.Create(&product)
	if err != nil {
		t.Error("Failed to add prepare product for test")
	}

	t.Run("Test for deletion", func(t *testing.T) {
		err = productsRepository.Delete(fetchedProduct)
		if err != nil {
			t.Errorf("Failed to delete product with id %d", product.Id)
		}

		t.Run("Try to fetch deleted product", func(t *testing.T) {
			fetchedProduct, err = productsRepository.FindById(product.Id)
			if err.Error() != ErrorProductNotFound {
				t.Errorf("Product with id %d was not deleted", product.Id)
			}
		})
	})

	t.Run("Try to delete non-existing product", func(t *testing.T) {
		fakeProduct := model.Product{
			Id:          1,
			Description: "Lauchzwiebeln",
			Ean:         5001819040871,
		}

		err = productsRepository.Delete(&fakeProduct)
		if err.Error() != ErrorProductDeletion {
			t.Errorf("Product with id %d was deleted", product.Id)
		}
	})
}
