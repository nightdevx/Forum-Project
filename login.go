package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
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

// Kullanıcının giriş bilgilerini doğrular
func authenticateUser(db *sql.DB, email, password string) (bool, error) {
	var storedPassword string
	query := "SELECT password FROM users WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	// Şifreyi karşılaştır (şifreleme olmadan)
	if password != storedPassword {
		return false, nil
	}
	return true, nil
}

// Rastgele oturum kimliği oluşturur
func generateSessionID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Kullanıcıyı giriş yapmış olarak ayarlar
func setSession(w http.ResponseWriter, email string, rememberMe bool) error {
	sessionID, err := generateSessionID()
	if err != nil {
		return err
	}

	var expiration time.Time
	if rememberMe {
		expiration = time.Now().Add(30 * 24 * time.Hour)
	} else {
		expiration = time.Now().Add(24 * time.Hour)
	}

	cookie := http.Cookie{
		Name:     "session_token",
		Value:    sessionID,
		Expires:  expiration,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	// Bu kısımda sessionID ve email'i bir session tablosuna kaydedebilirsiniz
	// veya başka bir oturum yönetim sistemi kullanabilirsiniz.

	// For demonstration purposes, we'll store the session in memory
	sessionStore[sessionID] = email

	return nil
}

// Store sessions in memory (for demonstration purposes)
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
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	authenticated, err := authenticateUser(db, email, password)
	if err != nil {
		http.Error(w, "Authentication error", http.StatusInternalServerError)
		return
	}

	if !authenticated {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	err = setSession(w, email, rememberMe)
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	// Ana sayfaya yönlendir
	http.Redirect(w, r, "/homepage", http.StatusSeeOther)
}

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	email, ok := sessionStore[cookie.Value]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Kullanıcının email adresini göster
	fmt.Fprintf(w, "Giriş başarılı! Hoş geldiniz, %s.", email)
}
