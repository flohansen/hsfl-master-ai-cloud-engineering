package products

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products/model"
	"reflect"
	"sort"
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

	productWithDublicateEan := model.Product{
		Id:          2,
		Description: "Dinkelnudeln",
		Ean:         4014819040771,
	}

	t.Run("Create product with success", func(t *testing.T) {
		_, err := demoRepository.Create(&product)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Check if products with duplicate ean can not be created", func(t *testing.T) {
		_, err := demoRepository.Create(&productWithDublicateEan)
		if err.Error() != ErrorEanAlreadyExists {
			t.Error(err)
		}
	})

	t.Run("Check for doublet", func(t *testing.T) {
		_, err := demoRepository.Create(&product)
		if err.Error() != ErrorProductAlreadyExists && err.Error() != ErrorEanAlreadyExists {
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
			sort.Slice(fetchedProducts, func(i, j int) bool {
				return fetchedProducts[i].Id < fetchedProducts[j].Id
			})
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

func TestDemoRepository_FindByEan(t *testing.T) {
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

	t.Run("Fetch product with existing ean", func(t *testing.T) {
		_, err := demoRepository.FindByEan(product.Ean)
		if err != nil {
			t.Errorf("Can't find expected product with ean %d", product.Ean)
		}

		t.Run("Is fetched product matching with added product?", func(t *testing.T) {
			fetchedProduct, _ := demoRepository.FindByEan(product.Ean)
			if !reflect.DeepEqual(product, *fetchedProduct) {
				t.Error("Fetched product does not match original product")
			}
		})
	})

	t.Run("Non-existing product test", func(t *testing.T) {
		_, err = demoRepository.FindByEan(42)
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

func TestDemoRepository_findNextAvailableID(t *testing.T) {
	type fields struct {
		products map[uint64]*model.Product
	}
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		{
			name: "Check if next available id is correct",
			fields: fields{products: map[uint64]*model.Product{
				1: {
					Id:          1,
					Description: "Strauchtomaten",
					Ean:         4014819040771,
				},
				2: {
					Id:          2,
					Description: "Lauchzwiebeln",
					Ean:         5001819040871,
				},
			}},
			want: 3,
		},
		{
			name: "Check if next available id is correct with gaps in map",
			fields: fields{products: map[uint64]*model.Product{
				1: {
					Id:          1,
					Description: "Strauchtomaten",
					Ean:         4014819040771,
				},
				3: {
					Id:          3,
					Description: "Lauchzwiebeln",
					Ean:         5001819040871,
				},
			}},
			want: 4,
		},
		{
			name:   "check if next available id is correct with empty map",
			fields: fields{products: make(map[uint64]*model.Product)},
			want:   1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &DemoRepository{
				products: tt.fields.products,
			}
			if got := repo.findNextAvailableID(); got != tt.want {
				t.Errorf("findNextAvailableID() = %v, want %v", got, tt.want)
			}
		})
	}
}
