package handler

import (
	"html/template"
	"log"
	"net/http"

	"url-shortener/middleware"
	"url-shortener/model"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.GetSessionStore().Get(r, "session")

	userIDRaw, ok := session.Values["user_id"]
	if !ok {
		log.Println("Session missing user_id")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var userID int
	switch v := userIDRaw.(type) {
	case int:
		userID = v
	case float64:
		userID = int(v)
	default:
		log.Println("Invalid session user_id type:", v)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	log.Println("Dashboard accessed by user_id:", userID)
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	links, err := model.GetLinksByUserID(userID)
	if err != nil {
		log.Println("Failed to fetch links:", err)
		http.Error(w, "Failed to fetch links", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		log.Println("Error parsing dashboard template:", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	successSlug := r.URL.Query().Get("success")

	// 🔧 Capture the execute error in its own variable
	execErr := tmpl.Execute(w, struct {
		Links   []model.ShortLink
		Error   string
		Success string
	}{
		Links:   links,
		Error:   "",
		Success: successSlug,
	})

	if execErr != nil {
		log.Println("Error rendering template:", execErr)
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}
