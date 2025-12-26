package main

import (
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var storedHash string
		err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedHash)

		// Compare password with hash
		if err == nil && bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)) == nil {
			// Set a simple session cookie
			cookie := &http.Cookie{
				Name:  "session",
				Value: username,
				Path:  "/",
			}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
	tmpl := template.Must(template.ParseFiles("templates/login.tmpl"))
	tmpl.Execute(w, nil)
}

func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		handler(w, r)
	}
}
