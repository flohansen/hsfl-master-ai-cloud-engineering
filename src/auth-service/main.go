package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	proto "github.com/Flo0807/hsfl-master-ai-cloud-engineering/lib/rpc/auth"
	"github.com/joho/godotenv"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/auth-service/api/http/handler"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/auth-service/api/http/router"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/auth-service/api/rpc"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/auth-service/auth"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/auth-service/config"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/auth-service/crypto"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/auth-service/user"
	"github.com/caarlos0/env/v10"
	"google.golang.org/grpc"
)

func GetenvInt(key string) int {
	value := os.Getenv(key)
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}

	return valueInt
}

func main() {
	godotenv.Load()

	cfg := config.Config{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("error while parsing enviroment variables: %s", err.Error())
	}

	userRepository, err := user.NewPsqlRepository(cfg.Database)

	if err != nil {
		log.Fatalf("error while creating user repository: %s", err.Error())
	}

	if err := userRepository.Migrate(); err != nil {
		log.Fatalf("could not migrate: %s", err.Error())
	}

	hasher := crypto.NewBcryptHasher()

	jwtTokenGenerator, err := auth.NewJwtTokenGenerator(cfg.JwtConfig)

	if err != nil {
		log.Fatalf("error while creating jwt token generator: %s", err.Error())
	}

	handler := router.NewRouter(
		handler.NewLoginHandler(userRepository, hasher, jwtTokenGenerator),
		handler.NewRegisterHandler(userRepository, hasher),
	)

	go func() {
		log.Printf("Starting HTTP server on port %s", cfg.HttpServerPort)

		addr := fmt.Sprintf("0.0.0.0:%s", cfg.HttpServerPort)
		if err := http.ListenAndServe(addr, handler); err != nil {
			log.Fatalf("error while listen and serve: %s", err.Error())
		}
	}()

	go func() {
		log.Printf("Starting gRPC server on port %s", cfg.GrpcServerPort)

		listener, err := net.Listen("tcp", ":"+cfg.GrpcServerPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		authServiceServer := rpc.NewAuthServiceServer(userRepository, jwtTokenGenerator)
		proto.RegisterAuthServiceServer(grpcServer, authServiceServer)

		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %s", err.Error())
		}
	}()

	select {}
}
