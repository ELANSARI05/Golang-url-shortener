package handler

import (
	"html/template"
	"log"
	"net/http"

	"url-shortener/model"
	"url-shortener/utils"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/signup.html")
		if err != nil {
			log.Println("Error loading signup template:", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	// Handle POST
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		tmpl, _ := template.ParseFiles("templates/signup.html")
		tmpl.Execute(w, map[string]string{"Error": "Username and password required"})
		return
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		log.Println("Hashing error:", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	err = model.CreateUser(username, hash)
	if err != nil {
		tmpl, _ := template.ParseFiles("templates/signup.html")
		tmpl.Execute(w, map[string]string{"Error": "Username already taken"})
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
