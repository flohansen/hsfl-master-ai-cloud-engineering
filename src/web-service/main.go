package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("SERVER_PORT")

	dir := http.Dir("./public")

	fs := http.FileServer(dir)

	mux := http.NewServeMux()

	mux.Handle("/", fs)

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
