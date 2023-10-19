package main

import (
	"flag"
	"fmt"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/api/router"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/transactions"
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
	port := flag.String("port", "8082", "The listening port")
	configPath := flag.String("config", "config.yaml", "The path to the configuration file")
	flag.Parse()

	config, err := LoadConfigFromFile(*configPath)
	if err != nil {
		log.Fatalf("could not load application configuration: %s", err.Error())
	}

	transactionRepository, err := transactions.NewPsqlRepository(config.Database)
	if err != nil {
		log.Fatalf("could not create user repository: %s", err.Error())
	}

	controller := transactions.NewDefaultController(transactionRepository)

	handler := router.New(controller)

	if err := transactionRepository.Migrate(); err != nil {
		log.Fatalf("could not migrate: %s", err.Error())
	}

	fmt.Println("Server started")

	addr := fmt.Sprintf("127.0.0.1:%s", *port)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
