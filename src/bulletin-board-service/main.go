package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/bulletin-board-service/api/handler"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/bulletin-board-service/api/router"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/bulletin-board-service/config"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/bulletin-board-service/models"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/bulletin-board-service/repository"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/bulletin-board-service/service"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load()

	cfg := config.Config{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("error while parsing enviroment variables: %s", err.Error())
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.Dsn()), &gorm.Config{})
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

	log.Printf("Starting HTTP server on port %s", cfg.HttpServerPort)

	addr := fmt.Sprintf("0.0.0.0:%s", cfg.HttpServerPort)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
