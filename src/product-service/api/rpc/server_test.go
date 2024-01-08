package rpc

import (
	"context"
	"github.com/stretchr/testify/assert"
	proto "hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/rpc/product"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products"
	"testing"
)

func TestNewProductServiceServer(t *testing.T) {
	var productRepo products.Repository = products.NewDemoRepository()
	var priceRepo prices.Repository = prices.NewDemoRepository()

	server := NewProductServiceServer(&productRepo, &priceRepo)

	assert.NotNil(t, server)
	assert.NotNil(t, server.productRepository)
	assert.NotNil(t, server.priceRepository, priceRepo)
}

func TestProductServiceServer_CreateProduct(t *testing.T) {
	productRepo := products.GenerateExampleDemoRepository()
	priceRepo := prices.GenerateExampleDemoRepository()
	server := NewProductServiceServer(&productRepo, &priceRepo)

	ctx := context.Background()
	req := &proto.CreateProductRequest{
		Product: &proto.Product{
			Id:          123,
			Description: "Test",
			Ean:         "123456789",
		},
	}

	resp, err := server.CreateProduct(ctx, req)
	if err != nil {
		t.Errorf("Error creating product: %v", err)
	}

	if resp.Product.Id != 123 {
		t.Errorf("Expected product ID '123', got '%d'", resp.Product.Id)
	}
}

func TestProductServiceServer_GetProduct(t *testing.T) {
	productRepo := products.GenerateExampleDemoRepository()
	priceRepo := prices.GenerateExampleDemoRepository()
	server := NewProductServiceServer(&productRepo, &priceRepo)

	ctx := context.Background()
	getReq := &proto.GetProductRequest{Id: 2}
	resp, err := server.GetProduct(ctx, getReq)
	if err != nil {
		t.Errorf("Error getting product: %v", err)
	}

	if resp.Product.Id != 2 {
		t.Errorf("Expected product ID '123', got '%d'", resp.Product.Id)
	}
}

func TestProductServiceServer_CreatePrice(t *testing.T) {

}

func TestProductServiceServer_DeletePrice(t *testing.T) {

}

func TestProductServiceServer_DeleteProduct(t *testing.T) {

}

func TestProductServiceServer_FindAllPrices(t *testing.T) {

}

func TestProductServiceServer_FindAllPricesFromUser(t *testing.T) {

}

func TestProductServiceServer_FindPrice(t *testing.T) {

}

func TestProductServiceServer_GetAllProducts(t *testing.T) {

}

func TestProductServiceServer_UpdatePrice(t *testing.T) {

}

func TestProductServiceServer_UpdateProduct(t *testing.T) {

}
