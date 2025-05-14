package main

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type ShortenRequest struct {
	URL        string `json:"url"`
	CustomSlug string `json:"custom_slug,omitempty"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

// ShortenHandler handles POST /shorten
func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ShortenRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &req)
	if err != nil || req.URL == "" {
		http.Error(w, "Invalid JSON or missing URL", http.StatusBadRequest)
		return
	}

	slug := req.CustomSlug
	if slug == "" {
		slug = generateSlug()
	} else {
		// Allow only letters, numbers, dashes, and underscores
		validSlug := regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)
		if !validSlug.MatchString(slug) {
			http.Error(w, "Invalid custom slug: only letters, numbers, hyphens, and underscores allowed", http.StatusBadRequest)
			return
		}
	}

	if err := saveURL(slug, req.URL); err != nil {
		http.Error(w, "Slug already exists", http.StatusConflict)
		return
	}

	resp := ShortenResponse{
		ShortURL: "http://localhost:8080/r/" + slug,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// RedirectHandler handles GET /{slug}
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/r/")
	if slug == "" {
		http.Error(w, "Slug required", http.StatusBadRequest)
		return
	}

	originalURL, found := getURL(slug)
	if !found {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

// GET /stats/{slug}
func statsHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/stats/")
	if slug == "" {
		http.Error(w, "Slug required", http.StatusBadRequest)
		return
	}

	stats, found := getStats(slug)
	if !found {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
func getAllHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(urlStore)
}
