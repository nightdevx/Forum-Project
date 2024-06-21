package main

import (
	"database/sql"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var tmpl = template.Must(template.ParseFiles("static/html/signup.html"))

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	username := r.FormValue("username")
	name := r.FormValue("name")
	surname := r.FormValue("surname")
	email := r.FormValue("email")
	password := r.FormValue("password")

	message := ""

	db, err := sql.Open("sqlite3", "database/forum.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var existingUsername, existingEmail string
	err = db.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&existingUsername)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = db.QueryRow("SELECT email FROM users WHERE email = ?", email).Scan(&existingEmail)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if existingUsername != "" {
		message += "Kullanıcı adı mevcut."
	}

	if existingEmail != "" {
		message += " E-posta adresi mevcut."
	}

	if message != "" {
		tmpl.Execute(w, struct{ Message string }{Message: message})
		return
	}

	// Kullanıcı adı ve e-posta adresi yoksa, yeni kayıt oluşturulur.
	stmt, err := db.Prepare("INSERT INTO users(username, name, surname, email, password) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, name, surname, email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Kayıt başarılı mesajı
	tmpl.Execute(w, struct{ Message string }{Message: "Kayıt başarılı!"})
}
