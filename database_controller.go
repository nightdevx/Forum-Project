package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func connectDatabase() error {
	var err error
	database, err = sql.Open("sqlite3", "database/forum.db")
	if err != nil {
		return err
	}
	return nil
}

func getUser(cookie *http.Cookie) (User, bool) {
	err := connectDatabase()
	if err != nil {
		return User{}, false
	}
	defer database.Close()

	var user User
	query := "SELECT username, email, name,surname,created_at,image FROM users WHERE id = ?"
	err = database.QueryRow(query, cookie.Value).Scan(&user.Username, &user.Email, &user.Name, &user.Surname, &user.CreationDate, &user.Image.ImageData)
	if err != nil {
		return User{}, false
	}
	return user, true
}

func getUserFullInfo(cookie *http.Cookie) (User, bool) {
	err := connectDatabase()
	if err != nil {
		fmt.Println(err)
		return User{}, false
	}

	var user User
	query := "SELECT username, email,biography, password,name,surname,image FROM users WHERE id = ?"
	err = database.QueryRow(query, cookie.Value).Scan(&user.Username, &user.Email, &user.Biography, &user.Password, &user.Name, &user.Surname, &user.Image.ImageData)
	if err != nil {
		fmt.Println(err)
		return User{}, false
	}
	fmt.Println(user)
	defer database.Close()

	return user, true
}

func updateUser(cookie *http.Cookie, user User) error {
	err := connectDatabase()
	if err != nil {
		fmt.Println(err)
		return err
	}
	// Hazırlık işlemi (Prepare statement)
	query, err := database.Prepare("UPDATE users SET username = ?, email = ?, biography = ?,name = ?, surname = ?, password = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer query.Close()

	// Parametreleri geçirerek sorguyu çalıştırma
	_, err = query.Exec(user.Username, user.Email, user.Biography, user.Name, user.Surname, user.Password, cookie.Value)
	if err != nil {
		return err
	}

	fmt.Println("Kullanıcı başarıyla güncellendi")
	return nil
}

func getPosts(userID int) ([]Post, error) {
	err := connectDatabase()
	if err != nil {
		return []Post{}, err
	}
	defer database.Close()

	query := `SELECT title, content,created_at,like_count,dislike_count FROM posts WHERE user_id = ?`
	rows, err := database.Query(query, userID)
	if err != nil {
		return []Post{}, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.PostTitle, &post.PostContent, &post.PostCreatedAt, &post.PostLikeCount, &post.PostDislikeCount); err != nil {
			return []Post{}, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return []Post{}, err
	}

	return posts, nil
}
func getAllPosts() []postData {
	err := connectDatabase()
	checkError(err)
	defer database.Close()
	rows, err := database.Query(`
		select users.username,users.name,users.surname, posts.title, posts.content, posts.created_at, posts.like_count, posts.dislike_count
		from posts
		join users on posts.user_id = users.id
		order by posts.created_at desc
	`)
	checkError(err)
	defer rows.Close()

	var posts []postData
	for rows.Next() {
		var tempPostData postData
		err = rows.Scan(&tempPostData.UserData.Username, &tempPostData.UserData.Name, &tempPostData.UserData.Surname, &tempPostData.PostData.PostTitle, &tempPostData.PostData.PostContent, &tempPostData.PostData.PostCreatedAt, &tempPostData.PostData.PostLikeCount, &tempPostData.PostData.PostDislikeCount)
		checkError(err)
		posts = append(posts, tempPostData)
	}
	return posts
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func insertPost(userID int, title, content string) {
	stmt, err := database.Prepare("insert into posts (user_id, title, content) values (?, ?, ?)")
	checkError(err)
	defer stmt.Close()

	_, err = stmt.Exec(userID, title, content)
	checkError(err)
}
