package main

import (
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/reverse-proxy/httpproxy"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
)

type RouteMapping struct {
	Path  string   `yaml:"path"`
	Hosts []string `yaml:"hosts"`
}

type ApplicationConfig struct {
	Mappings []RouteMapping `yaml:"mappings"`
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
	configFilePath := os.Getenv("CONFIG_FILE")
	if configFilePath == "" {
		configFilePath = "config.yaml"
	}

	config, err := LoadConfigFromFile(configFilePath)
	if err != nil {
		log.Fatalf("could not load application configuration: %s", err.Error())
	}

	proxy := httpproxy.New()

	for _, mapping := range config.Mappings {
		if err := proxy.Add(mapping.Path, mapping.Hosts); err != nil {
			log.Fatalf("Could not parse application config #{err.Error()}")
		}
	}

	log.Println("Server will be started!")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", proxy))
}
