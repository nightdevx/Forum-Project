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
	// Hashlenmiş şifreyi karşılaştır
	// err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	// if err != nil {
	// 	return false, nil
	// }
	// return true, nil

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
func setSession(w http.ResponseWriter, email string) error {
	sessionID, err := generateSessionID()
	if err != nil {
		return err
	}
	expiration := time.Now().Add(24 * time.Hour)

	cookie := http.Cookie{
		Name:     "session_token",
		Value:    sessionID,
		Expires:  expiration,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	// Bu kısımda sessionID ve email'i bir session tablosuna kaydedebilirsiniz
	// veya başka bir oturum yönetim sistemi kullanabilirsiniz.

	return nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.ServeFile(w, r, "./static/html/login.html")
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

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

	err = setSession(w, email)
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	// Giriş başarılı olduğunda mesajı yazdır
	fmt.Fprintf(w, "Giriş başarılı! Hoş geldiniz, %s.", email)

	// Ana sayfaya yönlendir
	http.Redirect(w, r, "/homepage", http.StatusSeeOther)
}
