package main

import (
	"log"
	"net/http"
)

func main() {
	err := loadFromFile()
	if err != nil {
		log.Fatalf("Failed to load data from file: %v", err)
	}

	// Serve index.html manually when accessing root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	// API endpoints
	http.HandleFunc("/shorten", shortenHandler) // for POSTing new links
	http.HandleFunc("/all", getAllHandler)      // for fetching all links as JSON
	http.HandleFunc("/stats/", statsHandler)    // optional
	http.HandleFunc("/r/", redirectHandler)     // handles short link redirection

	// Optional: serve other static assets (CSS/JS if needed later)
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server is running on http://localhost:8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
