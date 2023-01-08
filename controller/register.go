package controller

import (
	"forum/dbmanagement"
	utils "forum/helpers"
	"html/template"
	"net/http"

	"github.com/google/uuid"
)

func Register(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	if r.Method == "POST" {
		userName := r.FormValue("user_name")
		email := r.FormValue("email")
		password := utils.HashPassword(r.FormValue("password"))

		uuid := uuid.New()
		// log.Println(userName, email, password)

		dbmanagement.InsertUser(uuid.String(), userName, email, password, "user")
		dbmanagement.DisplayAllUsers()
	}

	tmpl.ExecuteTemplate(w, "register.html", nil)
}
