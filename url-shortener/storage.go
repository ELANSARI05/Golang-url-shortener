package main

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type URLRecord struct {
	OriginalURL string `json:"original_url"`
	Clicks      int    `json:"clicks"`
}

var (
	urlStore = make(map[string]URLRecord)
	mu       sync.RWMutex
)

// Save or update a new short URL
func saveURL(slug string, originalURL string) error {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := urlStore[slug]; exists {
		return errors.New("slug already exists")
	}
	urlStore[slug] = URLRecord{
		OriginalURL: originalURL,
		Clicks:      0,
	}
	return saveToFile()
}

// Retrieve a URL and increment click count
func getURL(slug string) (string, bool) {
	mu.Lock()
	defer mu.Unlock()

	record, exists := urlStore[slug]
	if !exists {
		return "", false
	}
	record.Clicks++
	urlStore[slug] = record
	saveToFile() // update click count
	return record.OriginalURL, true
}

// Get stats (clicks, etc.)
func getStats(slug string) (URLRecord, bool) {
	mu.RLock()
	defer mu.RUnlock()
	record, ok := urlStore[slug]
	return record, ok
}

// Save store to file
func saveToFile() error {
	file, err := os.Create("data.json")
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	return encoder.Encode(urlStore)
}

// Load store from file at startup
func loadFromFile() error {
	file, err := os.Open("data.json")
	if err != nil {
		return nil // it's okay if the file doesn't exist yet
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&urlStore)
}
