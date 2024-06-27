package main

import (
	"html/template"
	"net/http"
)

func PostPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		cookie, _ := r.Cookie("session_token")
		userID := cookie.Value
		commentContent := r.FormValue("comment")
		commentedPostID := r.FormValue("commentPostID")
		addCommentToDb(commentContent, commentedPostID, userID)

		// Extract query parameters from the original URL
		redirectURL := "/postpage?id=" + commentedPostID
		http.Redirect(w, r, redirectURL, http.StatusFound)

	} else if r.Method == http.MethodGet {
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
	}
}

