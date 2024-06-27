package main

import (
	"html/template"
	"net/http"
)

func PostPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		cookie, _ := r.Cookie("session_token")
		postID := r.FormValue("id")
		post, _ := getPostById(postID)
		var data struct {
			PostData   PostData
			IsLoggedIn bool
		}
		if cookie != nil {
			data = struct {
				PostData   PostData
				IsLoggedIn bool
			}{
				PostData:   post,
				IsLoggedIn: true,
			}
		} else {
			data = struct {
				PostData   PostData
				IsLoggedIn bool
			}{
				PostData:   post,
				IsLoggedIn: false,
			}
		}

		tmpl := template.Must(template.ParseFiles("static/html/postpage.html"))
		err := tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == "POST" {
		cookie, _ := r.Cookie("session_token")
		userID := cookie.Value
		commentContent := r.FormValue("comment")
		commentedPostID := r.FormValue("commentPostID")
		addCommentToDb(commentContent, commentedPostID, userID)
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
	}
}
