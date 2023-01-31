package main

import (
	"crypto/tls"
	"forum/controller"
	"forum/dbmanagement"
	"html/template"
	"log"
	"net/http"
	"os"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("static/*.html"))
	file, err := os.OpenFile("logfile.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		file, _ := os.Create("logfile.txt")
		defer file.Close()
	} else {
		defer file.Close()
	}
}

func protectGetRequests(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			controller.PageErrors(w, r, tmpl, "404")
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
	mux.HandleFunc("/authenticate", AuthenticateHandler)
	mux.HandleFunc("/logout", protectGetRequests(LogoutHandler))
	mux.HandleFunc("/register", protectGetRequests(RegisterHandler))
	mux.HandleFunc("/register_account", RegisterAccountHandler)

	// oauth handlers
	mux.HandleFunc("/google/login", protectGetRequests(GoogleLoginHandler))
	mux.HandleFunc("/google/callback", GoogleCallbackHandler)
	mux.HandleFunc("/github/login", protectGetRequests(GithubLoginHandler))
	mux.HandleFunc("/github/callback", GithubCallbackHandler)
	mux.HandleFunc("/facebook/login", protectGetRequests(FacebookLoginHandler))
	mux.HandleFunc("/facebook/callback", FacebookCallbackHandler)

	// forum handlers
	mux.HandleFunc("/forum", ForumHandler)
	mux.HandleFunc("/submitpost", SubmitPostHandler)
	mux.HandleFunc("/admin", AdminHandler)
	mux.HandleFunc("/user", UserHandler)
	mux.HandleFunc("/privacy_policy", PrivacyPolicyHandler)
	mux.HandleFunc("/error", ErrorHandler)

	// dbmanagement.DeleteUser("Yell Tro")
	// dbmanagement.CreateDatabaseWithTables()
	// dbmanagement.DeleteAllSessions()
	dbmanagement.DisplayAllUsers()
	log.Fatal(s.ListenAndServeTLS("", ""))
}
