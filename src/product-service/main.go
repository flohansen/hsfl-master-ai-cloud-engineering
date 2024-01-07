package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router/middleware/auth"
	proto "hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/rpc/product"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/rpc/user"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/api/http/router"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/api/rpc"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service/config"
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
	var configuration = loadConfiguration()

	var productRepository products.Repository = products.NewDemoRepository()
	var productsController products.Controller = products.NewCoalescingController(productRepository)
	createContentForProducts(productRepository)

	var priceRepository prices.Repository = prices.NewDemoRepository()
	var pricesController prices.Controller = prices.NewCoalescingController(priceRepository)
	createContentForPrices(priceRepository)

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go startHTTPServer(ctx, &wg, configuration, &productsController, &pricesController)

	wg.Add(1)
	go startGRPCServer(ctx, &wg, configuration, &productRepository, &priceRepository)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	<-stopChan
	cancel()

	wg.Wait()
}

func loadConfiguration() *config.ServiceConfiguration {
	godotenv.Load()

	serviceConfiguration := &config.ServiceConfiguration{}
	if err := env.Parse(serviceConfiguration); err != nil {
		log.Fatalf("couldn't parse configuration from environment: %s", err.Error())
	}

	return serviceConfiguration
}

func startHTTPServer(ctx context.Context, wg *sync.WaitGroup, configuration *config.ServiceConfiguration, productsController *products.Controller, pricesController *prices.Controller) {
	defer wg.Done()

	// Create client for user service for token validation
	userConn, err := grpc.Dial(configuration.GrpcUserServiceTarget,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer userConn.Close()
	grpcUserServiceClient := user.NewUserServiceClient(userConn)

	authMiddleware := auth.CreateAuthMiddleware(grpcUserServiceClient)
	handler := router.New(productsController, pricesController, authMiddleware)
	server := &http.Server{Addr: fmt.Sprintf(":%d", configuration.HttpPort), Handler: handler}

	go func() {
		log.Println("Starting HTTP server: ", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	<-ctx.Done()

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("HTTP Server Shutdown Failed:%v", err)
	}
}

func startGRPCServer(ctx context.Context, wg *sync.WaitGroup, configuration *config.ServiceConfiguration, productRepository *products.Repository, priceRepository *prices.Repository) {
	defer wg.Done()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", configuration.GrpcPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	productServiceServer := rpc.NewProductServiceServer(productRepository, priceRepository)
	proto.RegisterProductServiceServer(grpcServer, productServiceServer)
	proto.RegisterPriceServiceServer(grpcServer, productServiceServer)

	go func() {
		log.Println("Starting gRPC server: ", lis.Addr().String())
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
			Ean:         "4014819040771",
		},
		{
			Id:          2,
			Description: "Lauchzwiebeln",
			Ean:         "5001819040871",
		},
		{
			Id:          3,
			Description: "Mehl",
			Ean:         "5001819049871",
		},
	}

	for _, product := range productSlice {
		_, err := productRepository.Create(product)
		if err != nil {
			return
		}
	}
}
