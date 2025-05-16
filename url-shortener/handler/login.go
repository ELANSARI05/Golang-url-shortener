package handler

import (
	"html/template"
	"log"
	"net/http"

	"url-shortener/middleware"
	"url-shortener/model"
	"url-shortener/utils"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			log.Println("Error parsing login template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := model.GetUserByUsername(username)
	if err != nil || user == nil || !utils.CheckPasswordHash(password, user.PasswordHash) {
		log.Println("Login failed for:", username)
		tmpl, _ := template.ParseFiles("templates/login.html")
		tmpl.Execute(w, map[string]string{"Error": "Invalid username or password"})
		return
	}
	//Set session
	session, err := middleware.GetSessionStore().Get(r, "session")
	if err != nil {
		log.Println("Failed to get session:", err)
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	session.Values["user_id"] = user.ID
	if err := session.Save(r, w); err != nil {
		log.Println("Failed to save session:", err)
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	log.Println("Login success for:", username)
	log.Println("Session user_id:", session.Values["user_id"])

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.GetSessionStore().Get(r, "session")
	session.Options.MaxAge = -1
	session.Save(r, w)
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
