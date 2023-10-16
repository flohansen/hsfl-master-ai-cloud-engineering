package main

import (
	"log"
	"net/http"

	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router"
)

func main() {
	handler := router.New()

	if err := http.ListenAndServe(":3001", handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
