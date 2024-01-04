package rpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	proto "hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/rpc/product"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices"
	priceModel "hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices/model"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products"
	productsModel "hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products/model"
)

type ProductServiceServer struct {
	proto.UnimplementedProductServiceServer
	proto.UnimplementedPriceServiceServer
	productRepository *products.Repository
	priceRepository   *prices.Repository
}

func NewProductServiceServer(productRepository *products.Repository, priceRepository *prices.Repository) *ProductServiceServer {
	return &ProductServiceServer{
		productRepository: productRepository,
		priceRepository:   priceRepository,
	}
}

func (p *ProductServiceServer) CreateProduct(ctx context.Context, request *proto.CreateProductRequest) (*proto.CreateProductResponse, error) {
	requestProduct := request.GetProduct()

	var createdProduct, err = (*p.productRepository).Create(&productsModel.Product{
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

func (p *ProductServiceServer) GetProduct(ctx context.Context, request *proto.GetProductRequest) (*proto.GetProductResponse, error) {
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

func (p *ProductServiceServer) GetAllProducts(ctx context.Context, request *proto.GetAllProductsRequest) (*proto.GetAllProductsResponse, error) {
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

func (p *ProductServiceServer) UpdateProduct(ctx context.Context, request *proto.UpdateProductRequest) (*proto.UpdateProductResponse, error) {
	requestProduct := request.GetProduct()

	var updateProduct, err = (*p.productRepository).Update(&productsModel.Product{
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

func (p *ProductServiceServer) DeleteProduct(ctx context.Context, request *proto.DeleteProductRequest) (*proto.DeleteProductResponse, error) {
	requestProductId := request.Id

	var err = (*p.productRepository).Delete(&productsModel.Product{
		Id: requestProductId,
	})

	if err == nil {
		return &proto.DeleteProductResponse{}, nil
	} else {
		return nil, status.Error(codes.Internal, err.Error())
	}
}

func (p *ProductServiceServer) CreatePrice(ctx context.Context, request *proto.CreatePriceRequest) (*proto.CreatePriceResponse, error) {
	requestPrice := request.GetPrice()

	var createdPrice, err = (*p.priceRepository).Create(&priceModel.Price{
		UserId:    requestPrice.UserId,
		ProductId: requestPrice.ProductId,
		Price:     requestPrice.Price,
	})

	if err == nil {
		response := &proto.CreatePriceResponse{
			Price: &proto.Price{
				UserId:    createdPrice.UserId,
				ProductId: createdPrice.ProductId,
				Price:     createdPrice.Price,
			},
		}
		return response, nil
	} else {
		return nil, status.Error(codes.Internal, err.Error())
	}
}

func (p *ProductServiceServer) FindPrice(ctx context.Context, request *proto.FindPriceRequest) (*proto.FindPriceResponse, error) {
	var foundPrice, err = (*p.priceRepository).FindByIds(request.ProductId, request.UserId)

	if err == nil {
		response := &proto.FindPriceResponse{Price: &proto.Price{
			UserId:    foundPrice.UserId,
			ProductId: foundPrice.ProductId,
			Price:     foundPrice.Price,
		}}
		return response, nil
	} else {
		return nil, status.Error(codes.NotFound, err.Error())
	}
}

func (p *ProductServiceServer) FindAllPrices(ctx context.Context, request *proto.FindAllPricesRequest) (*proto.FindAllPricesResponse, error) {
	var foundPrices, err = (*p.priceRepository).FindAll()

	priceList := make([]*proto.Price, len(foundPrices))

	for i, price := range foundPrices {
		priceList[i] = &proto.Price{
			UserId:    price.UserId,
			ProductId: price.ProductId,
			Price:     price.Price,
		}
	}

	if err == nil {
		response := &proto.FindAllPricesResponse{
			Price: priceList,
		}
		return response, nil
	} else {
		return nil, status.Error(codes.Internal, err.Error())
	}
}

func (p *ProductServiceServer) FindAllPricesFromUser(ctx context.Context, request *proto.FindAllPricesFromUserRequest) (*proto.FindAllPricesFromUserResponse, error) {
	var foundPrices, err = (*p.priceRepository).FindAllByUser(request.UserId)

	priceList := make([]*proto.Price, len(foundPrices))

	for i, price := range foundPrices {
		priceList[i] = &proto.Price{
			UserId:    price.UserId,
			ProductId: price.ProductId,
			Price:     price.Price,
		}
	}

	if err == nil {
		response := &proto.FindAllPricesFromUserResponse{
			Price: priceList,
		}
		return response, nil
	} else {
		return nil, status.Error(codes.Internal, err.Error())
	}
}

func (p *ProductServiceServer) UpdatePrice(ctx context.Context, request *proto.UpdatePriceRequest) (*proto.UpdatePriceResponse, error) {
	requestPrice := request.GetPrice()

	var updatePrice, err = (*p.priceRepository).Update(&priceModel.Price{
		UserId:    requestPrice.UserId,
		ProductId: requestPrice.ProductId,
		Price:     requestPrice.Price,
	})

	if err == nil {
		response := &proto.UpdatePriceResponse{
			Price: &proto.Price{
				UserId:    updatePrice.UserId,
				ProductId: updatePrice.ProductId,
				Price:     updatePrice.Price,
			},
		}
		return response, nil
	} else {
		return nil, status.Error(codes.Internal, err.Error())
	}
}

func (p *ProductServiceServer) DeletePrice(ctx context.Context, request *proto.DeletePriceRequest) (*proto.DeletePriceResponse, error) {
	requestUserId := request.UserId
	requestProductId := request.ProductId

	var err = (*p.priceRepository).Delete(&priceModel.Price{
		UserId:    requestUserId,
		ProductId: requestProductId,
	})

	if err == nil {
		return &proto.DeletePriceResponse{}, nil
	} else {
		return nil, status.Error(codes.Internal, err.Error())
	}
}
