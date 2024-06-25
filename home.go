package main

import (
	"html/template"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_token")
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		userID, _ := strconv.Atoi(cookie.Value)
		insertPost(userID, title, content)
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
	currentUser, _ := getUser(cookie)
	homeData := struct {
		User  User
		Posts []postData
	}{
		User:  currentUser,
		Posts: getAllPosts(),
	}

	tmpl := template.Must(template.ParseFiles("static/html/homepage.html"))
	err := tmpl.Execute(w, homeData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
