package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func InitSessionStore() {
	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("SECRET environment variable not set")
	}

	store = sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		userID, ok := session.Values["user_id"]
		if !ok {
			log.Println("No session user_id, redirecting to login")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		log.Println("Authenticated user_id:", userID)
		next.ServeHTTP(w, r)
	})
}

func GetSessionStore() *sessions.CookieStore {
	return store
}
