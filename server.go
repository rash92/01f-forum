package main

import (
	"crypto/tls"
	auth "forum/authentication"
	"forum/controller"
	"forum/dbmanagement"
	"html/template"
	"log"
	"net/http"
	"os"
)

var tmpl *template.Template
var Tmpl *template.Template

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

	// index handlers
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			controller.PageErrors(w, r, tmpl, "404")
		} else {
			controller.AllPosts(w, r, tmpl)
		}
	})

	mux.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		tags := dbmanagement.SelectAllTags()
		tagexists := false
		var url string
		for _, v := range tags {
			if r.URL.Path == "/categories/"+v.TagName {
				url = v.TagName
				tagexists = true
			}
		}
		if !tagexists && r.URL.Path == "/" {
			controller.PageErrors(w, r, tmpl, "what ever you are looking for isn't here my dude")
		}
		if tagexists && r.URL.Path != "/" {
			controller.SubForum(w, r, tmpl, url)
		}
	})

	mux.HandleFunc("/posts/", func(w http.ResponseWriter, r *http.Request) {
		posts := dbmanagement.SelectAllPosts()
		postexists := false
		var url string
		for _, v := range posts {
			if r.URL.Path == "/posts/"+v.UUID {
				url = v.UUID
				postexists = true
			}
		}
		if !postexists && r.URL.Path == "/" {
			controller.PageErrors(w, r, tmpl, "what ever you are looking for isn't here my dude")
		}
		if postexists && r.URL.Path != "/" {
			controller.Post(w, r, tmpl, url)
		}
	})

	// authentication handlers
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		auth.Login(w, r, tmpl)
	})
	mux.HandleFunc("/authenticate", func(w http.ResponseWriter, r *http.Request) {
		auth.Authenticate(w, r, tmpl)
	})
	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		auth.Logout(w, r, tmpl)
	})
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		auth.Register(w, r, tmpl)
	})
	mux.HandleFunc("/register_account", func(w http.ResponseWriter, r *http.Request) {
		auth.RegisterAcount(w, r, tmpl)
	})

	// oauth handlers
	mux.HandleFunc("/google/login", func(w http.ResponseWriter, r *http.Request) {
		auth.GoogleLogin(w, r, tmpl)
	})
	mux.HandleFunc("/google/callback", func(w http.ResponseWriter, r *http.Request) {
		auth.GoogleCallback(w, r, tmpl)
	})
	mux.HandleFunc("/github/login", func(w http.ResponseWriter, r *http.Request) {
		auth.GithubLogin(w, r, tmpl)
	})
	mux.HandleFunc("/github/callback", func(w http.ResponseWriter, r *http.Request) {
		auth.GithubCallback(w, r, tmpl)
	})
	mux.HandleFunc("/facebook/login", func(w http.ResponseWriter, r *http.Request) {
		auth.FacebookLogin(w, r, tmpl)
	})
	mux.HandleFunc("/facebook/callback", func(w http.ResponseWriter, r *http.Request) {
		auth.FacebookCallback(w, r, tmpl)
	})

	// forum handlers
	mux.HandleFunc("/forum", func(w http.ResponseWriter, r *http.Request) {
		controller.AllPosts(w, r, tmpl)
	})
	mux.HandleFunc("/submitpost", func(w http.ResponseWriter, r *http.Request) {
		controller.SubmitPost(w, r, tmpl)
	})

	mux.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		controller.Admin(w, r, tmpl)
	})
	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		controller.User(w, r, tmpl)
	})
	mux.HandleFunc("/privacy_policy", func(w http.ResponseWriter, r *http.Request) {
		controller.PrivacyPolicy(w, r, tmpl)
	})
	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		controller.PageErrors(w, r, tmpl, "you're doing too much my dude")
	})

	// dbmanagement.DeleteUser("Yell Tro")
	// dbmanagement.CreateDatabaseWithTables()
	// dbmanagement.DeleteAllSessions()
	dbmanagement.DisplayAllUsers()
	log.Fatal(s.ListenAndServeTLS("", ""))
}
