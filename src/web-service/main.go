package main

import (
	"net/http"
)

type IndexPageViewModel struct {
	Title string
}

func main() {

	// Add a custom handler to set the Content-Type for JavaScript files.
	http.Handle("/", http.FileServer(http.Dir("frontend/static")))

	http.ListenAndServe(":3000", nil)
}
