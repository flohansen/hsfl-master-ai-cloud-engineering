package main

import (
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()
	router.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("dist"))))
	log.Fatal(http.ListenAndServe(":3000", router))
}
