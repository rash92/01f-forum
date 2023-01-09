package controller

import (
	"forum/dbmanagement"
	utils "forum/helpers"
	"html/template"
	"log"
	"net/http"
)

func RegisterAcount(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	if r.Method == "POST" {
		userName := r.FormValue("user_name")
		email := r.FormValue("email")
		password := utils.HashPassword(r.FormValue("password"))

		log.Println(userName, email, password)

		dbmanagement.InsertUser(userName, email, password, "user")
		dbmanagement.DisplayAllUsers()
	}

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}
