package controller

import (
	"forum/dbmanagement"
	"forum/security"
	"html/template"
	"log"
	"net/http"
)

/*
Registers a user with given details then redirects to log in page.  Password is hashed here.
*/
func RegisterAcount(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	if r.Method == "POST" {
		userName := r.FormValue("user_name")
		email := r.FormValue("email")
		password := security.HashPassword(r.FormValue("password"))

		log.Println(userName, email, password)

		dbmanagement.InsertUser(userName, email, password, "user")
		dbmanagement.DisplayAllUsers()
	}

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}
