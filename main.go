package main

import (
	"fmt"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	// http.HandleFunc("/", homePageHandler)
	http.Handle("/profile", http.HandlerFunc(profileHandler))
	http.Handle("/editProfile", http.HandlerFunc(editProfileHandler))
	// http.Handle("/sharePost", sessionMiddleware(http.HandlerFunc(sharePostHandler)))
	http.HandleFunc("/login", loginHandler)
	// http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/signup", SignupHandler)

	// Custom default handler to handle unknown routes
	defaultHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !knownRoutes(r.URL.Path) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			http.DefaultServeMux.ServeHTTP(w, r)
		}
	})

	fmt.Println("Server 8080 portu üzerinden başlatılıyor...")
	http.ListenAndServe(":8080", defaultHandler)
}

// Function to check if the URL path is a known route
func knownRoutes(path string) bool {
	knownPaths := []string{"/", "/profile", "/editProfile", "/sharePost", "/login", "/logout", "/signup", "/static/"}
	for _, p := range knownPaths {
		if path == p || (p == "/static/" && len(path) > len(p) && path[:len(p)] == p) {
			return true
		}
	}
	return false
}
