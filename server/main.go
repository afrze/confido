package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file if present
	_ = godotenv.Load(".env")

	// Get environment and port
	env := os.Getenv("APP_ENV")
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	if env == "production" {
		// Serve static files from the frontend build output
		staticDir := filepath.Join(".", "client", "dist")
		fs := http.FileServer(http.Dir(staticDir))

		// SPA fallback: serve index.html for non-file routes
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			requestedPath := filepath.Join(staticDir, filepath.Clean(r.URL.Path))
			if info, err := os.Stat(requestedPath); err == nil && !info.IsDir() {
				// File exists, serve it
				fs.ServeHTTP(w, r)
				return
			}
			// Fallback to index.html for SPA routing
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
		})

		log.Printf("Running in PRODUCTION mode: serving app on port %s", port)
	} else {
		// Development mode: no frontend served
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Frontend not served in development mode", http.StatusNotImplemented)
		})
	}

	log.Printf("Server listening on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
