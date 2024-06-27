package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_token")
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")

		err := r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error parsing form", http.StatusInternalServerError)
			return
		}

		postImage, _, err := r.FormFile("postPicture")
		var fileBytes []byte

		if err != nil {
			if err != http.ErrMissingFile {
				fmt.Println(err)
				log.Println("Error retrieving the profile picture")
			} else {
				fileBytes = []byte{} // No file uploaded, set to empty byte slice
			}
		} else {
			fileBytes, err = io.ReadAll(postImage)
			if err != nil {
				log.Println("Error reading file")
				return
			}
		}

		if len(title) >= 5 && len(content) >= 10 {
			userID, _ := strconv.Atoi(cookie.Value)
			insertPost(userID, title, content, fileBytes)
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		topPosts, err := getTopPosts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if cookie != nil {
			currentUser, _ := getUser(cookie)
			homeData := struct {
				User       User
				Posts      []PostData
				IsLoggedIn bool
				TopPosts   []PostData
			}{
				User:       currentUser,
				Posts:      getAllPosts(),
				IsLoggedIn: true,
				TopPosts:   topPosts,
			}

			tmpl := template.Must(template.ParseFiles("static/html/homepage.html"))
			err := tmpl.Execute(w, homeData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			homeData2 := struct {
				Posts      []PostData
				IsLoggedIn bool
				TopPosts   []PostData
			}{
				Posts:      getAllPosts(),
				IsLoggedIn: false,
				TopPosts:   topPosts,
			}

			tmpl := template.Must(template.ParseFiles("static/html/homepage.html"))
			err := tmpl.Execute(w, homeData2)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func getTopPosts() ([]PostData, error) {
	connectDatabase()
	rows, err := database.Query(`select users.username,users.name,users.surname, posts.id,posts.title, posts.content ,posts.created_at, posts.like_count, posts.dislike_count,posts.image
from posts join users on posts.user_id = users.id ORDER BY like_count DESC LIMIT 3`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []PostData
	for rows.Next() {
		var post PostData
		var image Image
		err = rows.Scan(&post.UserData.Username, &post.UserData.Name, &post.UserData.Surname, &post.PostData.PostID, &post.PostData.PostTitle, &post.PostData.PostContent, &post.PostData.PostCreatedAt, &post.PostData.PostLikeCount, &post.PostData.PostDislikeCount, &image.ImageData)
		post.PostData.PostImage = convertImg(image)
		checkError(err)
		posts = append(posts, post)
	}
	return posts, nil
}
