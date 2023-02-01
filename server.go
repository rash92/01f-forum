package main

import (
	"crypto/tls"
	"forum/controller"
	"forum/dbmanagement"
	"html/template"
	"log"
	"net/http"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("static/*.html"))
}

func protectGetRequests(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			controller.PageErrors(w, r, tmpl, 404, "Page Not Found")
		}
		h(w, r)
	}
}

func protectPostRequests(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			controller.PageErrors(w, r, tmpl, 404, "Page Not Found")
		}
		h(w, r)
	}
}

func main() {
	mux := http.NewServeMux()
	cert, _ := tls.LoadX509KeyPair("https/localhost.crt", "https/localhost.key")
	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}
	path := "./static"
	fs := http.FileServer(http.Dir(path))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// handlers
	mux.HandleFunc("/", protectGetRequests(IndexHandler))
	mux.HandleFunc("/posts", protectGetRequests(IndexHandler))
	mux.HandleFunc("/categories/", protectGetRequests(CategoriesHandler))
	mux.HandleFunc("/posts/", protectGetRequests(PostsHandler))

	// authentication handlers
	mux.HandleFunc("/login", protectGetRequests(LoginHandler))
	mux.HandleFunc("/authenticate", protectPostRequests(AuthenticateHandler))
	mux.HandleFunc("/logout", protectGetRequests(LogoutHandler))
	mux.HandleFunc("/register", protectGetRequests(RegisterHandler))
	mux.HandleFunc("/register_account", protectPostRequests(RegisterAccountHandler))

	// oauth handlers
	mux.HandleFunc("/google/login", protectPostRequests(GoogleLoginHandler))
	mux.HandleFunc("/google/callback", protectPostRequests(GoogleCallbackHandler))
	mux.HandleFunc("/github/login", protectPostRequests(GithubLoginHandler))
	mux.HandleFunc("/github/callback", protectPostRequests(GithubCallbackHandler))
	mux.HandleFunc("/facebook/login", protectGetRequests(FacebookLoginHandler))
	mux.HandleFunc("/facebook/callback", protectPostRequests(FacebookCallbackHandler))

	// forum handlers
	mux.HandleFunc("/forum", ForumHandler)
	mux.HandleFunc("/submitpost", protectPostRequests(SubmitPostHandler))
	mux.HandleFunc("/admin", protectPostRequests(AdminHandler))
	mux.HandleFunc("/user", protectPostRequests(UserHandler))
	mux.HandleFunc("/privacy_policy", protectGetRequests(PrivacyPolicyHandler))
	mux.HandleFunc("/error", protectGetRequests(ErrorHandler))

	// dbmanagement.DeleteUser("Yell Tro")
	// dbmanagement.CreateDatabaseWithTables()
	// dbmanagement.DeleteAllSessions()
	dbmanagement.ResetAllUserLoggedInStatus()
	// dbmanagement.DisplayAllUsers()
	log.Fatal(s.ListenAndServeTLS("", ""))
}
