package controller

import (
	"forum/dbmanagement"
	utils "forum/helpers"
	"html/template"
	"log"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	if r.Method == "POST" {
		userName := r.FormValue("user_name")
		// email := r.FormValue("email")
		// todo: ecnrypt password
		password := utils.HashPassword(r.FormValue("password"))

		log.Println(userName, password)
		user := dbmanagement.SelectUniqueUser(userName)

		if utils.CompareHash(user.Password, password) {
			log.Println("Password correct!")
			// session := user.CreateSession()
			http.Redirect(w, r, "/user", http.StatusFound)
			// /
		} else {
			log.Println("Incorrect Password!")
			http.Redirect(w, r, "/login", http.StatusFound)
		}

	}
	tmpl.ExecuteTemplate(w, "login.html", nil)
}
