package main

import (
	"flag"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/auth"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
)

type ApplicationConfig struct {
	Database database.PsqlConfig `yaml:"database"`
	Jwt      auth.JwtConfig      `yaml:"jwt"`
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

	tokenGenerator, err := auth.NewJwtTokenGenerator(config.Jwt)
	if err != nil {
		log.Fatalf("could not create JWT token generator: %s", err.Error())
	}
	claims := make(map[string]interface{})
	data, err := tokenGenerator.CreateToken(claims)
	if err != nil {
		log.Fatalf("could not create JWT token generator: %s", err.Error())
	}
	_ = data
	_ = tokenGenerator
	_ = port
	_ = config

	// Authentication stuff
	router.POST("/login", func(w http.ResponseWriter, r *http.Request) {

	})

	router.POST("/register", func(w http.ResponseWriter, r *http.Request) {

	})

	router.GET("/refresh-token", func(w http.ResponseWriter, r *http.Request) {

	})

	router.GET("/logout", func(w http.ResponseWriter, r *http.Request) {

	})

	// user-specific stuff
	router.GET("/users", func(w http.ResponseWriter, r *http.Request) {

	})

	router.GET("/users/:userId", func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId").(string)
		_ = userId
	})

	router.PUT("/users/:userId", func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId").(string)
		_ = userId
	})

	router.DELETE("/users/:userId", func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId").(string)
		_ = userId
	})

	router.GET("/users/:userId/books", func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId").(string)
		_ = userId
	})
}
