package main

import (
	"net/http"

	"url-shortener/handler"
	"url-shortener/middleware"
)

func RegisterRoutes() {
	// Auth routes
	http.HandleFunc("/signup", handler.SignupHandler)
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/logout", handler.LogoutHandler)

	//Protected routes (require authentication)
	http.Handle("/dashboard", middleware.Auth(http.HandlerFunc(handler.DashboardHandler)))
	http.Handle("/create", middleware.Auth(http.HandlerFunc(handler.CreateLinkHandler)))
	http.Handle("/delete", middleware.Auth(http.HandlerFunc(handler.DeleteLinkHandler)))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	//Shortlink redirect (last fallback route)
	http.HandleFunc("/", handler.RedirectHandler)

}
