package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/auth-service/api/handler"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/auth-service/api/router"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/auth-service/auth"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/auth-service/crypto"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/auth-service/database"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/auth-service/user"
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
	port := os.Getenv("PORT")

	psqlConfig := database.PsqlConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     GetenvInt("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Dbname:   os.Getenv("DB_NAME"),
	}

	userRepository, err := user.NewPsqlRepository(psqlConfig)

	if err != nil {
		log.Fatalf("error while creating user repository: %s", err.Error())
	}

	if err := userRepository.Migrate(); err != nil {
		log.Fatalf("could not migrate: %s", err.Error())
	}

	hasher := crypto.NewBcryptHasher()

	jwtConfig := auth.JwtConfig{
		SignKey: os.Getenv("JWT_SIGN_KEY"),
	}

	jwtTokenGenerator, err := auth.NewJwtTokenGenerator(jwtConfig)

	if err != nil {
		log.Fatalf("error while creating jwt token generator: %s", err.Error())
	}

	handler := router.NewRouter(
		handler.NewLoginHandler(userRepository, hasher, jwtTokenGenerator),
		handler.NewRegisterHandler(userRepository, hasher),
	)

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
