package main

import (
	"net/http"
	"strings"
)

type IndexPageViewModel struct {
	Title string
}

func main() {
	var outDir = "svelte/static"

	http.Handle("/", http.FileServer(http.Dir(outDir)))

	// Add a custom handler to set the Content-Type for JavaScript files.
	http.HandleFunc("/js/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
		}
		http.ServeFile(w, r, outDir+r.URL.Path)
	})

	http.ListenAndServe(":3000", nil)
}
