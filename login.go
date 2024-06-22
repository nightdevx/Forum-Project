package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Veritabanı bağlantısını açar
func openDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Kullanıcının giriş bilgilerini doğrular ve kullanıcı ID'sini döner
func authenticateUser(db *sql.DB, email, password string) (bool, int, error) {
	var storedPassword string
	var userID int
	query := "SELECT id, password FROM users WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&userID, &storedPassword)
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
	if r.Method != http.MethodPost {
		http.ServeFile(w, r, "./static/html/login.html")
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	rememberMe := r.FormValue("remember_me") == "on"

	db, err := openDatabase()
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		fmt.Println("Database connection error")
		return
	}
	defer db.Close()

	authenticated, userID, err := authenticateUser(db, email, password)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		fmt.Println("Authentication error")
		return
	}

	if !authenticated {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		fmt.Println("Invalid email or password")
		return
	}

	err = setSession(w, userID, email, rememberMe)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		fmt.Println("Session error")
		return
	}

	// Ana sayfaya yönlendir
	http.Redirect(w, r, "/homepage", http.StatusSeeOther)
}
