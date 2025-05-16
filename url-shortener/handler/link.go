package handler

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"url-shortener/db"
	"url-shortener/middleware"
	"url-shortener/model"
	"url-shortener/utils"
)

func CreateLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	session, _ := middleware.GetSessionStore().Get(r, "session")

	userIDRaw := session.Values["user_id"]
	var userID int
	switch v := userIDRaw.(type) {
	case int:
		userID = v
	case float64:
		userID = int(v)
	default:
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	originalURL := r.FormValue("url")
	slug := r.FormValue("slug")

	if originalURL == "" {
		renderDashboardWithError(w, userID, "URL is required")
		return
	}

	if slug == "" {
		slug = utils.GenerateSlug(6)
	} else {
		exists, err := model.SlugExists(slug)
		if err != nil {
			log.Println("Database error:", err)
			renderDashboardWithError(w, userID, "Database error while checking slug")
			return
		}
		if exists {
			renderDashboardWithError(w, userID, "Slug already taken")
			return
		}
	}

	err := model.CreateLink(userID, originalURL, slug)
	if err != nil {
		log.Println("CreateLink error:", err)
		renderDashboardWithError(w, userID, "Error saving short link")
		return
	}
	http.Redirect(w, r, "/dashboard?success="+slug, http.StatusSeeOther)

}
func renderDashboardWithError(w http.ResponseWriter, userID int, errorMsg string) {
	links, err := model.GetLinksByUserID(userID)
	if err != nil {
		log.Println("Error fetching links in error view:", err)
		http.Error(w, "Error loading dashboard", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		log.Println("Template error:", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, struct {
		Links   []model.ShortLink
		Error   string
		Success string
	}{
		Links:   links,
		Error:   errorMsg,
		Success: "",
	})

}
func DeleteLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	session, _ := middleware.GetSessionStore().Get(r, "session")
	userID := session.Values["user_id"].(int)

	linkID := r.FormValue("id")
	if linkID == "" {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	idInt, err := strconv.Atoi(linkID)
	if err != nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	model.DeleteLinkByID(userID, idInt)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Path[1:] // removes leading '/'
	if slug == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	url, err := model.GetOriginalURLBySlug(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	_, _ = db.DB.Exec(`
        UPDATE short_links SET click_count = click_count + 1 WHERE short_slug = ?
    `, slug)

	http.Redirect(w, r, url, http.StatusFound)
}
