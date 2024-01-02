//go:generate protoc --proto_path=../../lib/rpc --go_out=internal/proto --go_opt=paths=source_relative --go-grpc_out=internal/proto --go-grpc_opt=paths=source_relative ../../lib/rpc/product/product.proto
package main

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/api/http/router"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/api/rpc"
	proto "hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/internal/proto/product"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices"
	priceModel "hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/prices/model"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products"
	productModel "hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/products/model"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	var productRepository products.Repository = products.NewDemoRepository()
	var productsController products.Controller = products.NewCoalescingController(productRepository)
	createContentForProducts(productRepository)

	var priceRepository prices.Repository = prices.NewDemoRepository()
	var pricesController prices.Controller = prices.NewCoalescingController(priceRepository)
	createContentForPrices(priceRepository)

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go startHTTPServer(ctx, &wg, &productsController, &pricesController)

	wg.Add(1)
	go startGRPCServer(ctx, &wg, &productRepository, &priceRepository)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	<-stopChan
	cancel()

	wg.Wait()
}

func startHTTPServer(ctx context.Context, wg *sync.WaitGroup, productsController *products.Controller, pricesController *prices.Controller) {
	defer wg.Done()

	handler := router.New(productsController, pricesController)
	server := &http.Server{Addr: ":3003", Handler: handler}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	<-ctx.Done()

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("HTTP Server Shutdown Failed:%v", err)
	}
}

func startGRPCServer(ctx context.Context, wg *sync.WaitGroup, productRepository *products.Repository, priceRepository *prices.Repository) {
	defer wg.Done()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	productServiceServer := rpc.NewProductServiceServer(productRepository, priceRepository)
	proto.RegisterProductServiceServer(grpcServer, productServiceServer)
	proto.RegisterPriceServiceServer(grpcServer, productServiceServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	<-ctx.Done()
	grpcServer.GracefulStop()
}

func createContentForPrices(priceRepository prices.Repository) {
	pricesSlice := []*priceModel.Price{
		{
			UserId:    2,
			ProductId: 1,
			Price:     2.99,
		},
		{
			UserId:    2,
			ProductId: 2,
			Price:     5.99,
		},
		{
			UserId:    2,
			ProductId: 3,
			Price:     0.55,
		},
		{
			UserId:    1,
			ProductId: 3,
			Price:     0.55,
		},
	}

	for _, price := range pricesSlice {
		_, err := priceRepository.Create(price)
		if err != nil {
			return
		}
	}
}

func createContentForProducts(productRepository products.Repository) {
	productSlice := []*productModel.Product{
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
		{
			Id:          3,
			Description: "Mehl",
			Ean:         5001819049871,
		},
	}

	for _, product := range productSlice {
		_, err := productRepository.Create(product)
		if err != nil {
			return
		}
	}
}
