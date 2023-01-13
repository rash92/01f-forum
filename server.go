package main

import (
	"crypto/tls"
	"fmt"
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
	fmt.Println("got here")
	path := "static"
	fs := http.FileServer(http.Dir(path))
	// dbmanagement.UpdateUserPermissionFromName("admin", "admin")
	mux := http.NewServeMux()

	cert, _ := tls.LoadX509KeyPair("localhost.crt", "localhost.key")

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tags := dbmanagement.SelectAllTags()
		tagexists := false
		var url string
		for _, v := range tags {
			if r.URL.Path == "/"+v.TagName {
				url = v.TagName
				tagexists = true
			}
		}
		if !tagexists && r.URL.Path == "/" {
			controller.AllPosts(w, r, tmpl)
		}
		if tagexists && r.URL.Path != "/" {
			controller.SubForum(w, r, tmpl, url)
		}
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

	mux.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		controller.Admin(w, r, tmpl)
	})
	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		controller.User(w, r, tmpl)
	})

	dbmanagement.DisplayAllUsers()
	log.Fatal(s.ListenAndServeTLS("", ""))
}
