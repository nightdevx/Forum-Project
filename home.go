package main

import (
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		userID := 1
		insertPost(userID, title, content)
	}

	posts := getAllPosts()

	tmpl := template.Must(template.ParseFiles("static/html/homepage.html"))
	err := tmpl.Execute(w, posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
