package main

import (
	"flag"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/transactions"
	"gopkg.in/yaml.v3"
	"log"
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
	router := router.New()
	port := flag.String("port", "8080", "The listening port")
	configPath := flag.String("config", "config.yaml", "The path to the configuration file")
	flag.Parse()

	config, err := LoadConfigFromFile(*configPath)
	if err != nil {
		log.Fatalf("could not load application configuration: %s", err.Error())
	}

	userRepository, err := transactions.NewPsqlRepository(config.Database)
	if err != nil {
		log.Fatalf("could not create user repository: %s", err.Error())
	}

	if err := userRepository.Migrate(); err != nil {
		log.Fatalf("could not migrate: %s", err.Error())
	}

	_ = router
	_ = port

	/* handler := router.New(
		handler.NewRegisterHandler(userRepository, hasher),
		handler.NewLoginHandler(userRepository, hasher, tokenGenerator),
	)

	addr := fmt.Sprintf("0.0.0.0:%s", *port)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	} */
}
