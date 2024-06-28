package main

import (
	"fmt"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/home", homePageHandler)
	http.Handle("/profile", http.HandlerFunc(profileHandler))
	http.Handle("/editProfile", http.HandlerFunc(editProfileHandler))
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/signup", SignupHandler)
	http.HandleFunc("/profile/likepost", likePostHandler)
	http.HandleFunc("/profile/dislikepost", dislikePostHandler)
	http.HandleFunc("/home/likepost", likePostHandler)
	http.HandleFunc("/home/dislikepost", dislikePostHandler)
	http.HandleFunc("/postpage/likepost", dislikePostHandler)
	http.HandleFunc("/postpage/dislikepost", dislikePostHandler)
	http.HandleFunc("/postpage/likecomment", likeCommentHandler)
	http.HandleFunc("/postpage/dislikecomment", dislikeCommentHandler)
	http.HandleFunc("/postpage", PostPageHandler)
	http.HandleFunc("/discover", handleFilter)
	http.HandleFunc("/profile/likes", likesHandler)
	http.HandleFunc("/sifreyenileme", sifreyenilemeHandler)

	// Custom default handler to handle unknown routes
	defaultHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !knownRoutes(r.URL.Path) {
			http.Redirect(w, r, "/home", http.StatusSeeOther)
		} else {
			http.DefaultServeMux.ServeHTTP(w, r)
		}
	})

	fmt.Println("Server 8080 portu üzerinden başlatılıyor...")
	http.ListenAndServe(":8080", defaultHandler)
}

// Function to check if the URL path is a known route
func knownRoutes(path string) bool {
	knownPaths := []string{"/home", "/profile", "/editProfile", "/sharePost", "/login", "/logout", "/signup", "/profile/likepost", "/profile/dislikepost", "/home/likepost", "/discover", "/home/dislikepost", "/sifreyenileme", "/postpage", "/profile/likes", "/postpage/likecomment", "/postpage/dislikecomment", "/postpage/likepost", "/postpage/dislikepost", "/static/"}
	for _, p := range knownPaths {
		if path == p || (p == "/static/" && len(path) > len(p) && path[:len(p)] == p) {
			return true
		}
	}
	return false
}
