package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Kullanıcının giriş bilgilerini doğrular ve kullanıcı ID'sini döner
func authenticateUser(email, password string) (bool, int, error) {
	var storedPassword string
	var userID int
	query := "SELECT id, password FROM users WHERE email = ?"
	err := database.QueryRow(query, email).Scan(&userID, &storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, 0, nil
		}
		return false, 0, err
	}

	// Şifreyi karşılaştır (şifreleme olmadan)
	if password != storedPassword {
		return false, 0, nil
	}
	return true, userID, nil
}

// Kullanıcıyı giriş yapmış olarak ayarlar
func setSession(w http.ResponseWriter, userID int, email string, rememberMe bool) error {
	var expiration time.Time
	if rememberMe {
		expiration = time.Now().Add(30 * 24 * time.Hour)
	} else {
		expiration = time.Now().Add(24 * time.Hour)
	}

	cookie := http.Cookie{
		Name:     "session_token",
		Value:    strconv.Itoa(userID),
		Expires:  expiration,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	// Kullanıcı ID ve email'i bellekte sakla
	sessionStore[strconv.Itoa(userID)] = email

	return nil
}

// Bellekte oturumları sakla (demonstrasyon amacıyla)
var sessionStore = map[string]string{}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is already logged in by checking the session cookie
	cookie, err := r.Cookie("session_token")
	if err == nil {
		// If the cookie exists, check if the user is in the session store
		if _, ok := sessionStore[cookie.Value]; ok {
			// If the user is found in the session store, redirect to the homepage
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}
	}

	if r.Method != http.MethodPost {
		tmpl, _ := template.ParseFiles("./static/html/login.html")
		tmpl.Execute(w, nil)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	rememberMe := r.FormValue("remember_me") == "on"

	connectDatabase()

	authenticated, userID, err := authenticateUser(email, password)
	if err != nil {
		http.Error(w, "Authentication error", http.StatusInternalServerError)
		fmt.Println("Authentication error")
		return
	}

	if !authenticated {
		tmpl, _ := template.ParseFiles("./static/html/login.html")
		data := struct {
			ErrorMessage string
		}{
			ErrorMessage: "Geçersiz e-posta veya şifre",
		}
		w.WriteHeader(http.StatusUnauthorized) // Optional: set 401 Unauthorized status code
		tmpl.Execute(w, data)
		return
	}

	err = setSession(w, userID, email, rememberMe)
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		fmt.Println("Session error")
		return
	}

	// Redirect to the homepage
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
