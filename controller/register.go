package controller

import (
	"forum/dbmanagement"
	"html/template"
	"log"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	userName := r.FormValue("user_name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	log.Println(userName, email, password)
	dbmanagement.InsertUser("jkdgfhkj", userName, email, password, "user")
	dbmanagement.SelectUser()
}
