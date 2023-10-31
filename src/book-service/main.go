package main

import (
	"fmt"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/api/router"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

type ApplicationConfig struct {
	Database database.PsqlConfig `envPrefix:"POSTGRES_"`
	Port     uint16              `env:"PORT" envDefault:"8080"`
}

func main() {
	godotenv.Load()

	config := ApplicationConfig{}
	if err := env.Parse(&config); err != nil {
		log.Fatalf("Couldn't parse environment %s", err.Error())
	}

	bookRepository, err := books.NewPsqlBookRepository(config.Database)
	if err != nil {
		log.Fatalf("could not instanciate bookRepo: %s", err.Error())
	}
	chapterRepository, err := books.NewPsqlChapterRepository(config.Database)
	if err != nil {
		log.Fatalf("could not instanciate chapterRepo: %s", err.Error())
	}
	controller := books.NewDefaultController(bookRepository, chapterRepository)
	handler := router.New(controller)

	if err := bookRepository.Migrate(); err != nil {
		log.Fatalf("could not migrate: %s", err.Error())
	}
	if err := chapterRepository.Migrate(); err != nil {
		log.Fatalf("could not migrate: %s", err.Error())
	}

	log.Println("Server Started!")

	addr := fmt.Sprintf("0.0.0.0:%d", config.Port)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
