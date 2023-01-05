package controller

import (
	"forum/dbmanagement"
	"html/template"
	"log"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	if r.Method == "POST" {
		userName := r.FormValue("user_name")
		email := r.FormValue("email")
		password := r.FormValue("password")
		log.Println(userName, email, password)
		user := dbmanagement.SelectUniqueUser(userName)
		if CompareHash(user.Password, password) {
			http.Redirect(w, r, "http://localhost:8080/forum", http.StatusMovedPermanently)
			log.Println("Password correct!")
		} else {
			log.Println("Incorrent Password!")
		}
	}
	tmpl.ExecuteTemplate(w, "login.html", nil)

}
