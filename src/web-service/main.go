package main

import (
	"net/http"
)

type IndexPageViewModel struct {
	Title string
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("frontend/static")))
	http.ListenAndServe(":3000", nil)
}
