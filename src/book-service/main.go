package main

import (
	"flag"
	"fmt"
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
	port := flag.String("port", "8081", "The listening port")
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

	fmt.Println("Server started")

	addr := fmt.Sprintf("127.0.0.1:%s", *port)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
