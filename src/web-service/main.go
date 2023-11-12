package main

import (
	"net/http"
	"os"
)

func main() {
	fallbackFile := "index.html"
	rootDir := "frontend/static"

	http.Handle("/_app", http.FileServer(http.Dir("frontend/static/_app")))
	http.Handle("/", tryFilesHandler(rootDir, fallbackFile))

	http.ListenAndServe(":3000", nil)
}

func tryFilesHandler(rootDir string, fallbackFile string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Define the list of files to try, in the order specified
		tryFilesList := []string{r.URL.Path, r.URL.Path + ".html", ""}
		for _, file := range tryFilesList {
			filePath := rootDir + file
			if _, err := os.Stat(filePath); err == nil {
				http.ServeFile(w, r, filePath)
				return
			}
		}

		// Fallback to index.html if none of the files were found
		http.ServeFile(w, r, rootDir+fallbackFile)
	}
}
