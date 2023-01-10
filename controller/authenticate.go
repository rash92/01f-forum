package controller

import (
	"forum/dbmanagement"
	"forum/security"
	"forum/utils"
	"html/template"
	"log"
	"net/http"
)

/*
Displays the log in page.  If the username and password match an entry in the database then the user is redirected to the forum page, otherwise the user stays on the log in page.

Session Cookie is also set here.
*/
func Login(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	tmpl.ExecuteTemplate(w, "login.html", nil)
}

func Authenticate(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	// executed := false
	if r.Method == "POST" {
		userName := r.FormValue("user_name")
		password := r.FormValue("password")

		log.Println(userName, password)

		user := dbmanagement.SelectUserFromName(userName)

		if security.CompareHash(user.Password, password) {
			// log.Println("Password correct!")
			session, err := user.CreateSession()
			utils.HandleError("Cannot create user session err:", err)
			cookie := http.Cookie{
				Name:     "_cookie",
				Value:    session.UUID,
				HttpOnly: true,
			}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/forum", http.StatusMovedPermanently)
		} else {
			log.Println("Incorrect Password!")
			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		}

	}
	// if !executed {
	// 	tmpl.ExecuteTemplate(w, "login.html", nil)
	// }
}

// Logs the user out
func Logout(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	log.Println("logging out...")
	cookie, err := r.Cookie("_cookie")
	utils.HandleError("Failed to get cookie", err)

	if err != http.ErrNoCookie {
		session := cookie.Value
		// log.Println("session:", session)
		err := dbmanagement.DeleteSessionByUUID(session)
		utils.HandleError("Failed to get cookie", err)
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

/*
Displays the register page...
*/
func Register(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	tmpl.ExecuteTemplate(w, "register.html", nil)
}

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
