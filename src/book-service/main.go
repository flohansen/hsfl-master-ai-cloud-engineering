package main

import (
	"flag"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/api/router"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/database"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
)

type ApplicationConfig struct {
	Database database.PsqlConfig `yaml:"database"`
}

func LoadConfigFromFile(path string) (*ApplicationConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var config ApplicationConfig
	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	configPath := flag.String("config", "config.yaml", "The path to the configuration file")
	flag.Parse()

	config, err := LoadConfigFromFile(*configPath)
	if err != nil {
		log.Fatalf("could not load application configuration: %s", err.Error())
	}

	bookRepository, err := books.NewPsqlBookRepository(config.Database)
	chapterRepository, err := books.NewPsqlChapterRepository(config.Database)
	controller := books.NewDefaultController(bookRepository, chapterRepository)
	handler := router.New(controller)

	if err := bookRepository.Migrate(); err != nil {
		log.Fatalf("could not migrate: %s", err.Error())
	}

	if err := chapterRepository.Migrate(); err != nil {
		log.Fatalf("could not migrate: %s", err.Error())
	}

	if err := http.ListenAndServe(":3000", handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
