package main

import (
	"log"
	"net/http"
)

func main() {
	dir := http.Dir("./src/web-service/public")

	fs := http.FileServer(dir)

	mux := http.NewServeMux()

	mux.Handle("/", fs)
	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
