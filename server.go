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
		Addr:    ":8081",
		Handler: mux,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}
	path := "./static"
	fs := http.FileServer(http.Dir(path))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		controller.AllPosts(w, r, tmpl)
		// controller.Login(w, r, tmpl)
	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		controller.Login(w, r, tmpl)
	})
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		controller.Register(w, r, tmpl)
	})

	mux.HandleFunc("/register_account", func(w http.ResponseWriter, r *http.Request) {
		controller.RegisterAcount(w, r, tmpl)
	})

	mux.HandleFunc("/forum", func(w http.ResponseWriter, r *http.Request) {
		controller.AllPosts(w, r, tmpl)
	})
	mux.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		controller.Post(w, r, tmpl)
	})

	dbmanagement.DisplayAllUsers()
	log.Fatal(s.ListenAndServeTLS("", ""))
}
