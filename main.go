package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

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
	// mainRouter.Handle("/", frontendRouter)
	mainRouter.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			frontendRouter.ServeHTTP(w, r)
			return
		}

		// Handle 404 here
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 Not Found - No matching URL Patterns.\n 1. '/'\n 2. '/api/v1/*")
		log.Printf("[ ERR ] [%s] Route Not Found: %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)

		// fmt.Fprint(w, "")
	}))

	mainRouter.Handle("/api/v1/", http.StripPrefix("/api/v1", apiRouter))

	fmt.Println("Majiup running at PORT 8081")
	http.ListenAndServe(":8081", mainRouter)
}
