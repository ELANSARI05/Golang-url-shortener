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

	// Route to handle URL shortening
	http.HandleFunc("/shorten", shortenHandler)

	// Route to handle redirection
	http.HandleFunc("/", redirectHandler)

	http.HandleFunc("/stats/", statsHandler)

	log.Println("Server is running on http://localhost:8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
