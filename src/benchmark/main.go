package main

import (
	"flag"
	"log"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/benchmark/config"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/benchmark/loadtest"
)

func main() {
	configPath := flag.String("config", "./config.json", "Path to the configuration file")
	flag.Parse()

	config, err := config.LoadConfig(*configPath)

	if err != nil {
		log.Fatal("Error loading config file: ", err)
	}

	loadTest := loadtest.NewLoadTest(config)
	loadTest.Run()
}
