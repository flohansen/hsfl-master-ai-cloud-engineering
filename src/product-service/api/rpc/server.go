package rpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	proto "hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/internal/proto/product"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products/model"
)

type ProductServer struct {
	proto.UnimplementedProductServiceServer
	productRepository *products.Repository
	priceRepository   *prices.Repository
}

func NewProductServer(productRepository *products.Repository, priceRepository *prices.Repository) *ProductServer {
	return &ProductServer{
		productRepository: productRepository,
		priceRepository:   priceRepository,
	}
}

func (p *ProductServer) CreateProduct(ctx context.Context, request *proto.CreateProductRequest) (*proto.CreateProductResponse, error) {
	requestProduct := request.GetProduct()

	var createdProduct, err = (*p.productRepository).Create(&model.Product{
		Id:          requestProduct.GetId(),
		Description: requestProduct.GetDescription(),
		Ean:         requestProduct.GetEan(),
	})

	if err == nil {
		response := &proto.CreateProductResponse{
			Product: &proto.Product{
				Id:          createdProduct.Id,
				Description: createdProduct.Description,
				Ean:         createdProduct.Ean,
			},
		}
		return response, nil
	} else {
		return nil, status.Error(codes.Internal, err.Error())
	}
}

func (p *ProductServer) GetProduct(ctx context.Context, request *proto.GetProductRequest) (*proto.GetProductResponse, error) {
	var foundProduct, err = (*p.productRepository).FindById(request.GetId())

	if err == nil {
		response := &proto.GetProductResponse{
			Product: &proto.Product{
				Id:          foundProduct.Id,
				Description: foundProduct.Description,
				Ean:         foundProduct.Ean,
			},
		}
		return response, nil
	} else {
		return nil, status.Error(codes.NotFound, err.Error())
	}
}

func (p *ProductServer) GetAllProducts(ctx context.Context, request *proto.GetAllProductsRequest) (*proto.GetAllProductsResponse, error) {
	var foundProducts, err = (*p.productRepository).FindAll()

	productList := make([]*proto.Product, len(foundProducts))

	for i, product := range foundProducts {
		productList[i] = &proto.Product{
			Id:          product.Id,
			Description: product.Description,
			Ean:         product.Ean,
		}
	}

	if err == nil {
		response := &proto.GetAllProductsResponse{
			Products: productList,
		}
		return response, nil
	} else {
		return nil, status.Error(codes.Internal, err.Error())
	}
}

func (p *ProductServer) UpdateProduct(ctx context.Context, request *proto.UpdateProductRequest) (*proto.UpdateProductResponse, error) {
	requestProduct := request.GetProduct()

	var updateProduct, err = (*p.productRepository).Update(&model.Product{
		Id:          requestProduct.GetId(),
		Description: requestProduct.GetDescription(),
		Ean:         requestProduct.GetEan(),
	})

	if err == nil {
		response := &proto.UpdateProductResponse{
			Product: &proto.Product{
				Id:          updateProduct.Id,
				Description: updateProduct.Description,
				Ean:         updateProduct.Ean,
			},
		}
		return response, nil
	} else {
		return nil, status.Error(codes.Internal, err.Error())
	}
}

func (p *ProductServer) DeleteProduct(ctx context.Context, request *proto.DeleteProductRequest) (*proto.DeleteProductResponse, error) {
	requestProductId := request.Id

	var err = (*p.productRepository).Delete(&model.Product{
		Id: requestProductId,
	})

	if err == nil {
		return &proto.DeleteProductResponse{}, nil
	} else {
		return nil, status.Error(codes.Internal, err.Error())
	}
}
