package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func likesHandler(w http.ResponseWriter, r *http.Request) {
	// Kullanıcı kimliğini çerezden al
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Yetkisiz", http.StatusUnauthorized)
		log.Println("Cookie alınamadı:", err)
		return
	}

	kullaniciID := cookie.Value

	begendigiGonderiler, err := begendigiGonderileriGetir(kullaniciID)
	if err != nil {
		http.Error(w, "Bir şeyler yanlış gitti", http.StatusInternalServerError)
		log.Println("Gönderiler getirilemedi:", err)
		return
	}

	tmpl, err := template.ParseFiles("./static/html/likes.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Şablon dosyası yüklenemedi:", err)
		return
	}

	err = tmpl.Execute(w, begendigiGonderiler)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Şablon dosyası işlenemedi:", err)
	}
}

func begendigiGonderileriGetir(kullaniciID string) ([]PostData, error) {
	db, err := sql.Open("sqlite3", "database/forum.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
        SELECT posts.id, posts.title, posts.content, posts.image, posts.category, posts.like_count, posts.dislike_count, posts.created_at, users.username, users.name,users.surname
        FROM posts
        JOIN users ON posts.user_id = users.id
        JOIN likes ON posts.id = likes.post_id
        WHERE likes.user_id = ?
        ORDER BY likes.created_at DESC`
	rows, err := db.Query(query, kullaniciID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var postsData []PostData
	for rows.Next() {
		var postData PostData
		var image sql.NullString
		var category sql.NullString

		if err := rows.Scan(&postData.PostData.PostID, &postData.PostData.PostTitle, &postData.PostData.PostContent, &image, &category, &postData.PostData.PostLikeCount, &postData.PostData.PostDislikeCount, &postData.PostData.PostCreatedAt, &postData.UserData.Username, &postData.UserData.Name, &postData.UserData.Surname); err != nil {
			return nil, err
		}

		if image.Valid {
			imgData := Image{
				ImageData: []byte(image.String),
			}
			postData.PostData.PostImage = convertImg(imgData)
		}

		// Gönderiyi atan kullanıcının bilgilerini al
		postsData = append(postsData, postData)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return postsData, nil
}
