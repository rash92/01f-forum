package auth

import (
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type Data struct {
	ListOfData    []dbmanagement.Post
	Cookie        string
	UserInfo      dbmanagement.User
	TitleName     string
	IsCorrect     bool
	IsLoggedIn    bool
	RegisterError string
	TagsList      []dbmanagement.Tag
}

type OauthAccount struct {
	Name, Email string
}

// Displays the log in page.
func Login(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	data := Data{}
	data.TitleName = "Login"
	data.IsCorrect = true
	data.TagsList = dbmanagement.SelectAllTags()
	tmpl.ExecuteTemplate(w, "login.html", data)
}

/*
Authenticate user with credentials - If the username and password match an entry in the database then the user is redirected to the forum page,
otherwise the user stays on the log in page. Session Cookie is also set here.
*/
func Authenticate(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	if r.Method == "POST" {
		userName := r.FormValue("user_name")
		password := r.FormValue("password")

		user, err := dbmanagement.SelectUserFromName(userName)
		utils.HandleError("unable to get user error:", err)

		if CompareHash(user.Password, password) && user.IsLoggedIn == 0 {
			err := CreateUserSession(w, r, user)
			utils.HandleError("Failed to create session in authenticate", err)
			dbmanagement.UpdateUserLoggedInStatus(user.UUID, 1)
			log.Println(user.IsLoggedIn)
			http.Redirect(w, r, "/forum", http.StatusSeeOther)
		} else {
			if user.IsLoggedIn != 0 {
				log.Println("Already Logged In!")
				data := Data{}
				data.TitleName = "Login"
				data.IsCorrect = true
				data.IsLoggedIn = true
				data.TagsList = dbmanagement.SelectAllTags()
				tmpl.ExecuteTemplate(w, "login.html", data)
			} else {
				log.Println("Incorrect Password!")
				data := Data{}
				data.TitleName = "Login"
				data.IsCorrect = false
				data.TagsList = dbmanagement.SelectAllTags()
				tmpl.ExecuteTemplate(w, "login.html", data)
			}
		}
	}
}

// Logs user out
func Logout(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	log.Println("logging out...")
	cookie, err := r.Cookie("session")
	log.Println("Current Cookie: ", cookie)
	utils.HandleError("Failed to get cookie", err)
	session := cookie.Value
	user, _ := dbmanagement.SelectUserFromSession(session)
	dbmanagement.UpdateUserLoggedInStatus(user.UUID, 0)
	log.Println(user.IsLoggedIn)
	if err != http.ErrNoCookie {
		err := dbmanagement.DeleteSessionByUUID(session)
		utils.HandleError("Failed to get cookie", err)
	}
	clearcookie := http.Cookie{
		Name:     "session",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
	}

	http.SetCookie(w, &clearcookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Displays the register page
func Register(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	data := Data{}
	data.TagsList = dbmanagement.SelectAllTags()
	tmpl.ExecuteTemplate(w, "register.html", data)
}

// Registers a user with given details then redirects to log in page.  Password is hashed here.
func RegisterAcount(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	if r.Method == "POST" {
		userName := r.FormValue("user_name")
		email := r.FormValue("email")
		password := HashPassword(r.FormValue("password"))
		_, err := dbmanagement.InsertUser(userName, email, password, "user", 0)
		data := Data{}
		if err != nil {
			data.RegisterError = strings.Split(err.Error(), ".")[1]
			data.TagsList = dbmanagement.SelectAllTags()
			tmpl.ExecuteTemplate(w, "register.html", data)
		}
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
