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

func main() {
	mux := http.NewServeMux()
	cert, _ := tls.LoadX509KeyPair("localhost.crt", "localhost.key")
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

	// index handlers
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		controller.AllPosts(w, r, tmpl)
	})

	// authentication handlers
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		controller.Login(w, r, tmpl)
	})
	mux.HandleFunc("/authenticate", func(w http.ResponseWriter, r *http.Request) {
		controller.Authenticate(w, r, tmpl)
	})
	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		controller.Logout(w, r, tmpl)
	})
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		controller.Register(w, r, tmpl)
	})
	mux.HandleFunc("/register_account", func(w http.ResponseWriter, r *http.Request) {
		controller.RegisterAcount(w, r, tmpl)
	})

	// google authentication handlers
	mux.HandleFunc("/google/login", func(w http.ResponseWriter, r *http.Request) {
		controller.GoogleLogin(w, r, tmpl)
	})
	mux.HandleFunc("/google/callback", func(w http.ResponseWriter, r *http.Request) {
		controller.GoogleCallback(w, r, tmpl)
	})

	// forum handlers
	mux.HandleFunc("/forum", func(w http.ResponseWriter, r *http.Request) {
		controller.UsersPosts(w, r, tmpl)
	})

	dbmanagement.DisplayAllUsers()
	log.Fatal(s.ListenAndServeTLS("", ""))
}
