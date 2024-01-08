package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/load-balancer/config"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cfg := config.Config{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("error while parsing enviroment variables: %s", err.Error())
	}

	dir := http.Dir("./public")

	fs := http.FileServer(dir)

	mux := http.NewServeMux()

	mux.Handle("/", fs)

	log.Printf("Starting HTTP server on port %s", cfg.HttpServerPort)

	addr := fmt.Sprintf("0.0.0.0:%s", cfg.HttpServerPort)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
