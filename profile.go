package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func profileHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("ID")
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	userData, hasUser := getUser(cookie)
	userID, _ := strconv.Atoi(cookie.Value)
	postData, _ := getPosts(userID)
	if !hasUser {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	profileData := struct {
		User  User
		Posts []Post
		Img   string
	}{
		User:  userData,
		Posts: postData,
		Img:   convertImg(userData.Image),
	}
	tmpl, err := template.ParseFiles("static/html/profile.html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, profileData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func convertImg(img Image) string {
	imageBase64 := base64.StdEncoding.EncodeToString(img.ImageData)
	return imageBase64
}

func editProfileHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("ID")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	userData, hasUser := getUserFullInfo(cookie)
	if !hasUser {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	editData := struct {
		User    User
		Img     string
		Changes ChangeMessage
	}{
		User: userData,
		Img:  convertImg(userData.Image),
		Changes: ChangeMessage{
			Message:   "",
			IsChanged: false,
		},
	}

	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("static/html/editProfile.html")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, editData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == "POST" {
		userData.Username = r.FormValue("username")
		userData.Email = r.FormValue("email")
		userData.Name = r.FormValue("name")
		userData.Surname = r.FormValue("surname")
		newPassword := r.FormValue("newPassword")
		oldPassword := r.FormValue("oldPassword")
		if oldPassword != userData.Password && oldPassword != "" {
			editData = struct {
				User    User
				Img     string
				Changes ChangeMessage
			}{
				User: userData,
				Img:  convertImg(userData.Image),
				Changes: ChangeMessage{
					Message:   "Eski şifrenizi yanlış girdiniz!",
					IsChanged: true,
				},
			}
			tmpl, err := template.ParseFiles("static/html/editProfile.html")
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			err = tmpl.Execute(w, editData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		} else if oldPassword == newPassword && oldPassword != "" && newPassword != "" {
			editData = struct {
				User    User
				Img     string
				Changes ChangeMessage
			}{
				User: userData,
				Img:  convertImg(userData.Image),
				Changes: ChangeMessage{
					Message:   "Yeni şifreniz eskisiyle aynı olamaz!",
					IsChanged: true,
				},
			}
			tmpl, err := template.ParseFiles("static/html/editProfile.html")
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			err = tmpl.Execute(w, editData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		} else if newPassword != "" && oldPassword == "" {
			editData = struct {
				User    User
				Img     string
				Changes ChangeMessage
			}{
				User: userData,
				Img:  convertImg(userData.Image),
				Changes: ChangeMessage{
					Message:   "Şifrenizi girmediniz!",
					IsChanged: true,
				},
			}
			tmpl, err := template.ParseFiles("static/html/editProfile.html")
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			err = tmpl.Execute(w, editData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		userData.Password = newPassword
		err := updateUser(cookie, userData)
		if err != nil {
			fmt.Println("Error in update user", err)
			return
		} else {
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
		}
	}
}
