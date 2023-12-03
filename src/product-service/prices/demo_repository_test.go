package prices

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices/model"
	"reflect"
	"testing"
)

func TestNewDemoRepository(t *testing.T) {
	t.Run("Demo repository correct initialized", func(t *testing.T) {
		got := NewDemoRepository()
		if len(got.prices) != 0 {
			t.Errorf("NewDemoRepository().prices has unexpected length, want 0")
		}
	})
}

func TestDemoRepository_Create(t *testing.T) {
	demoRepository := NewDemoRepository()

	price := model.Price{
		UserId:    1,
		ProductId: 1,
		Price:     2.99,
	}

	t.Run("Create price with success", func(t *testing.T) {
		_, err := demoRepository.Create(&price)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Check for doublet", func(t *testing.T) {
		_, err := demoRepository.Create(&price)
		if err.Error() != ErrorPriceAlreadyExists {
			t.Error(err)
		}
	})
}

func TestDemoRepository_FindAll(t *testing.T) {
	demoRepository := NewDemoRepository()

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

	for _, price := range prices {
		_, err := demoRepository.Create(price)
		if err != nil {
			t.Error("Failed to add prepared price for test")
		}
	}

	t.Run("Fetch all prices", func(t *testing.T) {
		fetchedPrices, err := demoRepository.FindAll()
		if err != nil {
			t.Error("Can't fetch prices")
		}

		if len(fetchedPrices) != len(prices) {
			t.Errorf("Unexpected price count. Expected %d, got %d", len(prices), len(fetchedPrices))
		}
	})

	priceTest := []struct {
		name string
		want *model.Price
	}{
		{
			name: "first",
			want: prices[0],
		},
		{
			name: "second",
			want: prices[1],
		},
	}

	for _, tt := range priceTest {
		t.Run("Is fetched price matching with "+tt.name+" added price?", func(t *testing.T) {
			fetchedPrices, _ := demoRepository.FindAll()
			found := false

			for _, fetchedPrice := range fetchedPrices {
				if reflect.DeepEqual(tt.want, fetchedPrice) {
					found = true
					break
				}
			}

			if !found {
				t.Error("Fetched price does not match original price")
			}
		})
	}
}

func TestDemoRepository_FindAllByUser(t *testing.T) {
	demoRepository := NewDemoRepository()

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

	for _, price := range prices {
		_, err := demoRepository.Create(price)
		if err != nil {
			t.Fatalf("Failed to add prepared price for test: %v", err)
		}
	}

	t.Run("Fetch prices for a specific user", func(t *testing.T) {
		fetchedPrices, err := demoRepository.FindAllByUser(2)
		if err != nil {
			t.Fatalf("Error fetching prices: %v", err)
		}

		if len(fetchedPrices) != 1 {
			t.Errorf("Unexpected price count. Expected 1, got %d", len(fetchedPrices))
		}

		if !reflect.DeepEqual(prices[1], fetchedPrices[0]) {
			t.Error("Fetched price does not match the expected price")
		}
	})

	t.Run("Verify fetched prices match added prices", func(t *testing.T) {
		fetchedPrices, err := demoRepository.FindAll()
		if err != nil {
			t.Fatalf("Error fetching prices: %v", err)
		}

		for _, tt := range prices {
			found := false

			for _, fetchedPrice := range fetchedPrices {
				if reflect.DeepEqual(tt, fetchedPrice) {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("Fetched price for user id %d does not match the original price", tt.UserId)
			}
		}
	})
}

func TestDemoRepository_FindByIds(t *testing.T) {
	demoRepository := NewDemoRepository()

	price := model.Price{
		UserId:    1,
		ProductId: 1,
		Price:     2.99,
	}

	_, err := demoRepository.Create(&price)
	if err != nil {
		t.Fatal("Failed to add prepare price for test")
	}

	t.Run("Fetch price with existing ids", func(t *testing.T) {
		_, err := demoRepository.FindByIds(price.ProductId, price.UserId)
		if err != nil {
			t.Errorf("Can't find expected price with product id %d and user id %d", price.ProductId, price.UserId)
		}

		t.Run("Is fetched price matching with added price?", func(t *testing.T) {
			fetchedPrice, _ := demoRepository.FindByIds(price.ProductId, price.UserId)
			if !reflect.DeepEqual(price, *fetchedPrice) {
				t.Error("Fetched price does not match original price")
			}
		})
	})

	t.Run("Non-existing price test", func(t *testing.T) {
		_, err = demoRepository.FindByIds(42, 42)
		if err.Error() != ErrorPriceNotFound {
			t.Error(err)
		}
	})
}

func TestDemoRepository_Update(t *testing.T) {
	demoRepository := NewDemoRepository()

	price := model.Price{
		UserId:    1,
		ProductId: 1,
		Price:     2.99,
	}

	fetchedPrice, err := demoRepository.Create(&price)
	if err != nil {
		t.Error("Failed to add prepare price for test")
	}

	t.Run("Check if updated price object has updated price", func(t *testing.T) {
		price := model.Price{
			UserId:    1,
			ProductId: 1,
			Price:     3.99,
		}

		updatedPrice, err := demoRepository.Update(&price)
		if err != nil {
			t.Error(err.Error())
		}

		if fetchedPrice.Price != updatedPrice.Price {
			t.Errorf("Failed to update price. Got %f, want %f.",
				fetchedPrice.Price, updatedPrice.Price)
		}
	})
}

func TestDemoRepository_Delete(t *testing.T) {
	// Prepare test
	demoRepository := NewDemoRepository()

	price := model.Price{
		UserId:    1,
		ProductId: 1,
		Price:     2.99,
	}

	fetchedPrice, err := demoRepository.Create(&price)
	if err != nil {
		t.Error("Failed to add prepare price for test")
	}

	t.Run("Test for deletion", func(t *testing.T) {
		err = demoRepository.Delete(fetchedPrice)
		if err != nil {
			t.Errorf("Failed to delete price with product id %d and user id %d", price.ProductId, price.UserId)
		}

		t.Run("Try to fetch deleted price", func(t *testing.T) {
			fetchedPrice, err = demoRepository.FindByIds(price.ProductId, price.UserId)
			if err.Error() != ErrorPriceNotFound {
				t.Errorf("Price with with product id %d and user id %d was not deleted", price.ProductId, price.UserId)
			}
		})
	})

	t.Run("Try to delete non-existing price", func(t *testing.T) {
		fakePrice := model.Price{
			UserId:    2,
			ProductId: 2,
			Price:     5.99,
		}

		err = demoRepository.Delete(&fakePrice)
		if err.Error() != ErrorPriceDeletion {
			t.Errorf("Price with product id %d and user id %d was deleted", price.ProductId, price.UserId)
		}
	})
}
