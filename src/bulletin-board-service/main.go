package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/bulletin-board-service/api/handler"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/bulletin-board-service/api/router"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/bulletin-board-service/database"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/bulletin-board-service/models"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/bulletin-board-service/repository"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/bulletin-board-service/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	port := os.Getenv("SERVER_PORT")

	psqlConfig := database.PsqlConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     GetenvInt("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Dbname:   os.Getenv("DB_NAME"),
	}

	db, err := gorm.Open(postgres.Open(psqlConfig.Dsn()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	err = db.AutoMigrate(&models.Post{})
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	postRepo := repository.NewPostPsqlRepository(db)
	postService := service.NewPostService(postRepo)
	postHandler := handler.NewPostHandler(postService)

	healthHandler := handler.NewHealthHandler()

	r := router.NewRouter(healthHandler, postHandler)

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
