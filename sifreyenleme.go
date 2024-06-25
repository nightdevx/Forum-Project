package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Şifre sıfırlama işleyicisi
func sifreyenilemeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("static/html/sifreyenileme.html")
		if err != nil {
			http.Error(w, "Şablon dosyası yüklenemedi", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	} else if r.Method == http.MethodPost {
		email := r.FormValue("email")
		tckimlikLast3 := r.FormValue("tckimlik_last3")

		connectDatabase()
		var userID int
		var storedTCKimlik string
		err := database.QueryRow("SELECT id, tckimlik FROM users WHERE email = ?", email).Scan(&userID, &storedTCKimlik)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Bu e-posta ile kayıtlı kullanıcı bulunamadı", http.StatusNotFound)
				return
			}
			http.Error(w, "Veritabanı sorgu hatası", http.StatusInternalServerError)
			return
		}

		// Güvenlik sorusunun cevabını kontrol et
		if strings.ToLower(tckimlikLast3) != strings.ToLower(storedTCKimlik[len(storedTCKimlik)-3:]) {
			http.Error(w, "Güvenlik sorusunun cevabı yanlış", http.StatusUnauthorized)
			return
		}

		// Yeni şifre oluştur
		newPassword, err := generateRandomPassword()
		if err != nil {
			http.Error(w, "Yeni şifre oluşturulamadı", http.StatusInternalServerError)
			return
		}

		// Yeni şifreyi veritabanında güncelle
		err = updatePassword(userID, newPassword)
		if err != nil {
			http.Error(w, "Şifre güncellenemedi", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Yeni şifreniz başarıyla güncellendi. %s", newPassword)
	}
}

// Rastgele 6 karakter uzunluğunda şifre oluşturma
func generateRandomPassword() (string, error) {
	bytes := make([]byte, 3)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Veritabanında şifre güncelleme
func updatePassword(userID int, newPassword string) error {
	connectDatabase()
	_, err := database.Exec("UPDATE users SET password = ? WHERE id = ?", newPassword, userID)
	if err != nil {
		return err
	}

	return nil
}
