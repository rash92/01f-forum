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
		password := r.FormValue("password")

		log.Println(userName, password)

		user := dbmanagement.SelectUserFromName(userName)

		if utils.CompareHash(user.Password, password) {
			log.Println("Password correct!")
			session, err := user.CreateSession()
			utils.HandleError("Cannot create user session err:", err)
			cookie := http.Cookie{
				Name:     "_cookie",
				Value:    session.UUID,
				HttpOnly: true,
			}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/forum", http.StatusFound)
		} else {
			log.Println("Incorrect Password!")
			http.Redirect(w, r, "/login", http.StatusFound)
		}

	}
	tmpl.ExecuteTemplate(w, "login.html", nil)

}
