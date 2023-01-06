package main

import (
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
	path := "static"
	fs := http.FileServer(http.Dir(path))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		controller.Home(w, r, tmpl)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		controller.Login(w, r, tmpl)
	})
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		controller.Register(w, r, tmpl)
	})

	// dbmanagement.CreateDatabase()
	// dbmanagement.ModifyDatabase()
	dbmanagement.CreateDatabaseWithTables()
	dbmanagement.DisplayAllUsers()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
