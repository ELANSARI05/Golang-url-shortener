package main

import (
	"log"
	"net/http"

	"url-shortener/db"
	"url-shortener/middleware"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from config.env
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading config.env file:", err)
	}

	//Initialize session store after loading env
	middleware.InitSessionStore()

	db.Init()

	RegisterRoutes()

	log.Println("Server started at :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
