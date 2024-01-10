package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	proto "hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/rpc/user"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/api/http/handler"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/api/http/middleware"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/api/http/router"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/api/rpc"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/auth"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/config"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/crypto"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user/model"
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

	var usersRepository user.Repository = user.NewDemoRepository()
	usersRepository = createMockRepository(usersRepository)
	var usersController user.Controller = user.NewDefaultController(usersRepository)

	var tokenGenerator = createTokenGenerator(configuration)

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go startHTTPServer(ctx, &wg, configuration, &usersController, &usersRepository, tokenGenerator)

	wg.Add(1)
	go startGRPCServer(ctx, &wg, configuration, &usersRepository, tokenGenerator)

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

func startHTTPServer(ctx context.Context, wg *sync.WaitGroup, configuration *config.ServiceConfiguration, usersController *user.Controller, usersRepository *user.Repository, tokenGenerator auth.TokenGenerator) {
	defer wg.Done()

	var loginHandler = createLoginHandler(*usersRepository, tokenGenerator)
	var registerHandler = createRegisterHandler(*usersRepository)

	authMiddleware := middleware.CreateLocalAuthMiddleware(usersRepository, tokenGenerator)
	handler := router.New(loginHandler, registerHandler, usersController, authMiddleware)
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

func startGRPCServer(ctx context.Context, wg *sync.WaitGroup, configuration *config.ServiceConfiguration, usersRepository *user.Repository, tokenGenerator auth.TokenGenerator) {
	defer wg.Done()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", configuration.GrpcPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userServiceServer := rpc.NewUserServiceServer(usersRepository, tokenGenerator)
	proto.RegisterUserServiceServer(grpcServer, userServiceServer)

	go func() {
		log.Println("Starting gRPC server: ", lis.Addr().String())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	<-ctx.Done()
	grpcServer.GracefulStop()
}

func createTokenGenerator(configuration *config.ServiceConfiguration) auth.TokenGenerator {
	tokenGenerator, err := auth.NewJwtTokenGenerator(auth.JwtConfig{PrivateKey: configuration.JwtConfig.PrivateKey})
	if err != nil {
		panic(fmt.Sprintf("Can't generate token generator: %v", err))
	}
	return tokenGenerator
}

func createLoginHandler(userRepository user.Repository, tokenGenerator auth.TokenGenerator) *handler.LoginHandler {
	return handler.NewLoginHandler(userRepository,
		crypto.NewBcryptHasher(), tokenGenerator)
}

func createRegisterHandler(userRepository user.Repository) *handler.RegisterHandler {
	return handler.NewRegisterHandler(userRepository,
		crypto.NewBcryptHasher())
}

func createMockRepository(userRepository user.Repository) user.Repository {
	userSlice := createDemoUserSlice()
	for _, newUser := range userSlice {
		_, _ = userRepository.Create(newUser)
	}

	return userRepository
}

func createDemoUserSlice() []*model.User {
	bcryptHasher := crypto.NewBcryptHasher()
	hashedPassword, _ := bcryptHasher.Hash([]byte("12345"))

	return []*model.User{
		{
			Id:       1,
			Email:    "ada.lovelace@gmail.com",
			Password: hashedPassword,
			Name:     "Ada Lovelace",
			Role:     model.Customer,
		},
		{
			Id:       2,
			Email:    "info-aldi@gmail.com",
			Password: hashedPassword,
			Name:     "Aldi",
			Role:     model.Merchant,
		},
		{
			Id:       3,
			Email:    "info-edeka@gmail.com",
			Password: hashedPassword,
			Name:     "Edeka",
			Role:     model.Merchant,
		},
	}
}
