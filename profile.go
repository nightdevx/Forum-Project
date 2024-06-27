package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func profileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		userData, hasUser := getUser(cookie)
		if !hasUser {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userID := cookie.Value
		postData, _ := getPosts(userID)

		profileData := struct {
			User      User
			Posts     []Post
			MostLiked []PostData
		}{
			User:      userData,
			Posts:     postData,
			MostLiked: getUsersTopPosts(userID),
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
}

func editProfileHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println("No session token found")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	userData, hasUser := getUserFullInfo(cookie)
	if !hasUser {
		fmt.Println("User not found")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	editData := struct {
		User User
	}{
		User: userData,
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
		userData.Biography = r.FormValue("bio")
		newPassword := r.FormValue("newPassword")
		oldPassword := r.FormValue("oldPassword")
		if userData.Username == "" || userData.Email == "" || userData.Name == "" {
			http.Redirect(w, r, "/editprofile", http.StatusSeeOther)
			return
		}
		if newPassword != "" {
			if oldPassword != userData.Password && oldPassword != "" {
				http.Redirect(w, r, "/editprofile", http.StatusSeeOther)
				return
			} else if oldPassword == newPassword && oldPassword != "" && newPassword != "" {
				http.Redirect(w, r, "/editprofile", http.StatusSeeOther)
				return
			} else if newPassword != "" && oldPassword == "" {
				http.Redirect(w, r, "/editprofile", http.StatusSeeOther)
				return
			}
			userData.Password = newPassword
		}

		err := r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error parsing form", http.StatusInternalServerError)
			return
		}

		profileImage, _, err := r.FormFile("profilePicture")
		if err != nil && err != http.ErrMissingFile {
			fmt.Println(err)
			log.Println("Error retrieving the profile picture")
		}

		bannerImage, _, err := r.FormFile("bannerPicture")
		if err != nil && err != http.ErrMissingFile {
			fmt.Println(err)
			log.Println("Error retrieving the banner picture")
		}

		if profileImage != nil {
			defer profileImage.Close()
			uploadFile(r, profileImage, "image")
		}

		if bannerImage != nil {
			defer bannerImage.Close()
			uploadFile(r, bannerImage, "banner")
		}

		err = updateUser(cookie, userData)
		if err != nil {
			fmt.Println("Error updating user", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else {
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
		}
	}
}

func likePostHandler(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_token")
	if cookie == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	postID := r.FormValue("id")
	userID := cookie.Value
	isLiked := likesController(userID, postID)
	isDisliked := dislikesController(userID, postID)
	if isLiked {
		deleteLikedPost(userID, postID)
		decreaseLikeCount(postID)
	} else if !isLiked && !isDisliked {
		insertLikedPost(userID, postID)
		err := increaseLikeCount(postID)
		if err != nil {
			fmt.Println("Error liking post", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	if strings.HasPrefix(r.URL.Path, "/profile") {
		http.Redirect(w, r, "/profile", http.StatusFound)
	} else if strings.HasPrefix(r.URL.Path, "/home") {
		http.Redirect(w, r, "/home", http.StatusFound)
	}
}

func insertLikedPost(userID, postID string) {
	connectDatabase()
	stmt, err := database.Prepare("insert into likes (user_id, post_id) values (?, ?)")
	checkError(err)
	defer stmt.Close()

	_, err = stmt.Exec(userID, postID)
	checkError(err)
}

func likesController(userID, postID string) bool {
	connectDatabase()
	query := "SELECT COUNT(*) FROM likes WHERE user_id = ? AND post_id = ?"
	var count int
	err := database.QueryRow(query, userID, postID).Scan(&count)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return false
	}

	// Sonuçları kontrol et
	if count > 0 {
		fmt.Println("Record exists")
		return true
	} else {
		fmt.Println("Record does not exist")
		return false
	}
}

func deleteLikedPost(userID, postID string) {
	connectDatabase()
	query := "DELETE FROM likes WHERE user_id = ? AND post_id = ?"

	result, err := database.Exec(query, userID, postID)
	if err != nil {
		fmt.Println("Error executing delete query:", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Error getting rows affected:", err)
		return
	}

	if rowsAffected > 0 {
		fmt.Println("Like successfully removed")
	} else {
		fmt.Println("No like found to remove")
	}
}

func dislikePostHandler(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_token")
	if cookie == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	postID := r.FormValue("id")
	userID := cookie.Value
	isDisliked := dislikesController(userID, postID)
	isLiked := likesController(userID, postID)
	if isDisliked {
		deleteDislikedPost(userID, postID)
		decreaseDislikeCount(postID)
	} else if !isDisliked && !isLiked {
		insertDislikedPost(userID, postID)
		err := increaseDislikeCount(postID)
		if err != nil {
			fmt.Println("Error liking post", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	if strings.HasPrefix(r.URL.Path, "/profile") {
		http.Redirect(w, r, "/profile", http.StatusFound)
	} else if strings.HasPrefix(r.URL.Path, "/home") {
		http.Redirect(w, r, "/home", http.StatusFound)
	}
}

func dislikesController(userID, postID string) bool {
	connectDatabase()
	query := "SELECT COUNT(*) FROM dislikes WHERE user_id = ? AND post_id = ?"
	var count int
	err := database.QueryRow(query, userID, postID).Scan(&count)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return false
	}

	// Sonuçları kontrol et
	if count > 0 {
		fmt.Println("Record exists")
		return true
	} else {
		fmt.Println("Record does not exist")
		return false
	}
}

func insertDislikedPost(userID, postID string) {
	connectDatabase()
	stmt, err := database.Prepare("insert into dislikes (user_id, post_id) values (?, ?)")
	checkError(err)
	defer stmt.Close()

	_, err = stmt.Exec(userID, postID)
	checkError(err)
}

func deleteDislikedPost(userID, postID string) {
	connectDatabase()
	query := "DELETE FROM dislikes WHERE user_id = ? AND post_id = ?"

	result, err := database.Exec(query, userID, postID)
	if err != nil {
		fmt.Println("Error executing delete query:", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Error getting rows affected:", err)
		return
	}

	if rowsAffected > 0 {
		fmt.Println("Like successfully removed")
	} else {
		fmt.Println("No like found to remove")
	}
}

func uploadFile(r *http.Request, file io.Reader, target string) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error reading file")
		return
	}
	cookie, _ := r.Cookie("session_token")
	saveImageToDB(fileBytes, cookie.Value, target)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Çerez oluşturun ve süresini geçmiş bir zamana ayarlayın
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour), // Geçmiş bir zamana ayarlayın
		MaxAge:   -1,                             // Hemen silmek için MaxAge'yi -1 yapın
		HttpOnly: true,
	}

	// Çerezi HTTP yanıtına ekleyin
	http.SetCookie(w, cookie)

	// Kullanıcıyı anasayfaya yönlendirin veya başka bir işlem yapın
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func convertImg(img Image) string {
	imageBase64 := base64.StdEncoding.EncodeToString(img.ImageData)
	return imageBase64
}
