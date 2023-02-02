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

}

func protectGetRequests(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			controller.PageErrors(w, r, tmpl, 400, "Bad Request")
			return
		}
		h(w, r)

	}
}

func protectPostRequests(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			controller.PageErrors(w, r, tmpl, 400, "Bad Request")
			return
		}
		h(w, r)
	}
}

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--reset" {
		dbmanagement.CreateDatabaseWithTables()
	}

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
	// mux.HandleFunc("/posts", protectGetRequests(IndexHandler))
	mux.HandleFunc("/categories/", CategoriesHandler)
	mux.HandleFunc("/posts/", PostsHandler)

	// authentication handlers
	mux.HandleFunc("/login", protectGetRequests(LoginHandler))
	mux.HandleFunc("/authenticate", protectPostRequests(AuthenticateHandler))
	mux.HandleFunc("/logout", protectGetRequests(LogoutHandler))
	mux.HandleFunc("/register", protectGetRequests(RegisterHandler))
	mux.HandleFunc("/register_account", protectPostRequests(RegisterAccountHandler))

	// oauth handlers
	mux.HandleFunc("/google/login", GoogleLoginHandler)
	mux.HandleFunc("/google/callback", GoogleCallbackHandler)
	mux.HandleFunc("/github/login", GithubLoginHandler)
	mux.HandleFunc("/github/callback", GithubCallbackHandler)
	mux.HandleFunc("/facebook/login", FacebookLoginHandler)
	mux.HandleFunc("/facebook/callback", FacebookCallbackHandler)

	// forum handlers
	mux.HandleFunc("/forum", ForumHandler)
	mux.HandleFunc("/submitpost", SubmitPostHandler)
	mux.HandleFunc("/admin", AdminHandler)
	mux.HandleFunc("/user", UserHandler)
	mux.HandleFunc("/privacy_policy", PrivacyPolicyHandler)
	mux.HandleFunc("/error", ErrorHandler)

	// dbmanagement.DeleteUser("Yell Tro")
	dbmanagement.DeleteAllSessions()
	dbmanagement.ResetAllUserLoggedInStatus()
	dbmanagement.ResetAllTokens()
	// dbmanagement.DisplayAllUsers()
	log.Fatal(s.ListenAndServeTLS("", ""))
}
