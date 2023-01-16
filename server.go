package main

import (
	"crypto/tls"
	auth "forum/authentication"
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

// func HandleRoute(route string, wr func(w http.ResponseWriter, r *http.Request, tmpl *template.Template)) {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
// 		wr(w, r, tmpl)
// 	})
// }

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

	// github authentication handlers
	mux.HandleFunc("/github/login", func(w http.ResponseWriter, r *http.Request) {
		auth.GithubLogin(w, r, tmpl)
	})

	mux.HandleFunc("/github/callback", func(w http.ResponseWriter, r *http.Request) {
		auth.GithubCallback(w, r, tmpl)
	})

	// forum handlers
	mux.HandleFunc("/forum", func(w http.ResponseWriter, r *http.Request) {
		controller.AllPosts(w, r, tmpl)
	})
	mux.HandleFunc("/submitpost", func(w http.ResponseWriter, r *http.Request) {
		controller.SubmitPost(w, r, tmpl)
	})
	mux.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		controller.Post(w, r, tmpl)
	})

	mux.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		controller.Admin(w, r, tmpl)
	})
	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		controller.User(w, r, tmpl)
	})

	// dbmanagement.DeleteUser("rhemtro")
	// dbmanagement.DeleteAllSessions()
	dbmanagement.DisplayAllUsers()
	log.Fatal(s.ListenAndServeTLS("", ""))
}
