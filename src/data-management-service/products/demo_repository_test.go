package products

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/data-management-service/products/model"
	"reflect"
	"testing"
)

func TestProductsRepository_Create(t *testing.T) {
	// Prepare test
	demoRepository := NewDemoRepository()

	product := model.Product{
		Id:          1,
		Description: "Strauchtomaten",
		Ean:         4014819040771,
	}

	// Create product with success
	_, err := demoRepository.Create(&product)
	if err != nil {
		t.Error(err)
	}

	// Check for doublet
	_, err = demoRepository.Create(&product)
	if err.Error() != "product already exists" {
		t.Error(err)
	}
}

func TestProductsRepository_FindAll(t *testing.T) {
	// Prepare test
	demoRepository := NewDemoRepository()

	product1 := model.Product{
		Id:          1,
		Description: "Strauchtomaten",
		Ean:         4014819040771,
	}

	product2 := model.Product{
		Id:          2,
		Description: "Lauchzwiebeln",
		Ean:         5001819040871,
	}

	_, err := demoRepository.Create(&product1)
	if err != nil {
		t.Error("Failed to add prepare product for test")
	}

	_, err = demoRepository.Create(&product2)
	if err != nil {
		t.Error("Failed to add prepare product for test")
	}

	// Fetch all products
	fetchedProducts, err := demoRepository.FindAll()
	if err != nil {
		t.Error("Can't fetch products")
	}

	if len(fetchedProducts) != 2 {
		t.Errorf("Unexpected product count. Expected 2, got %d", len(fetchedProducts))
	}

	// Is fetched product matching with added product?
	if !reflect.DeepEqual(product1, *fetchedProducts[0]) {
		t.Error("Fetched product does not match original product")
	}

	// Is fetched product matching with added product?
	if !reflect.DeepEqual(product2, *fetchedProducts[1]) {
		t.Error("Fetched product does not match original product")
	}
}

func TestProductsRepository_FindById(t *testing.T) {
	// Prepare test
	demoRepository := NewDemoRepository()

	product := model.Product{
		Id:          1,
		Description: "Strauchtomaten",
		Ean:         4014819040771,
	}

	_, err := demoRepository.Create(&product)
	if err != nil {
		t.Error("Failed to add prepare product for test")
	}

	// Fetch product with existing id
	fetchedProduct, err := demoRepository.FindById(product.Id)
	if err != nil {
		t.Errorf("Can't find expected product with id %d", product.Id)
	}

	// Is fetched product matching with added product?
	if !reflect.DeepEqual(product, *fetchedProduct) {
		t.Error("Fetched product does not match original product")
	}

	// Non-existing product test
	_, err = demoRepository.FindById(42)
	if err.Error() != "product could not be found" {
		t.Error(err)
	}
}

func TestProductsRepository_Update(t *testing.T) {
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

	fetchedProduct.Description = "Wittenseer Mineralwasser"
	updatedProduct, err := demoRepository.Update(fetchedProduct)

	// Check if returned product has the updated description
	if fetchedProduct.Description != updatedProduct.Description {
		t.Error("Failed to update product")
	}
}

func TestProductsRepository_Delete(t *testing.T) {
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

	// Test for deletion
	err = productsRepository.Delete(fetchedProduct)
	if err != nil {
		t.Errorf("Failed to delete product with id %d", product.Id)
	}

	// Fetch product with existing id
	fetchedProduct, err = productsRepository.FindById(product.Id)
	if err.Error() != "product could not be found" {
		t.Errorf("Product with id %d was not deleted", product.Id)
	}

	// Try to delete non-existing product
	fakeProduct := model.Product{
		Id:          1,
		Description: "Lauchzwiebeln",
		Ean:         5001819040871,
	}

	// Test for deletion
	err = productsRepository.Delete(&fakeProduct)
	if err.Error() != "product could not be deleted" {
		t.Errorf("Product with id %d was deleted", product.Id)
	}
}
