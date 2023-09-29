package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/JosephMusya/majiup-backend/api"
	"github.com/gorilla/mux"
)

func main() {

	apiRouter := mux.NewRouter()

	api.ApiServe(apiRouter)

	// Create a new Gorilla Mux router
	frontendRouter := mux.NewRouter()

	appDir := "serve/dist"

	// Define a custom handler to serve JavaScript and CSS files with correct MIME types
	customHandler := func(w http.ResponseWriter, r *http.Request) {
		// Determine the file extension to set the correct Content-Type
		filePath := r.URL.Path
		ext := strings.ToLower(filepath.Ext(filePath))
		contentType := ""
		switch ext {
		case ".js":
			contentType = "application/javascript"
		case ".css":
			contentType = "text/css"
		}

		// Set the Content-Type header
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}

		// Serve the file
		http.ServeFile(w, r, filepath.Join(appDir, filePath))
	}

	// Register custom handler for JavaScript and CSS files
	frontendRouter.PathPrefix("/assets/").Handler(http.HandlerFunc(customHandler))

	// Serve the main HTML page
	frontendRouter.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmlFile := filepath.Join(appDir, "index.html")
		http.ServeFile(w, r, htmlFile)
	})

	mainRouter := http.NewServeMux()
	mainRouter.Handle("/", frontendRouter)
	mainRouter.Handle("/api/v1/", http.StripPrefix("/api/v1", apiRouter))

	fmt.Println("Majiup running at PORT 8081")
	http.ListenAndServe(":8081", mainRouter)
}
