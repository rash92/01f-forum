package controller

import (
	"fmt"
	"forum/dbmanagement"
	"html/template"
	"net/http"
)

type AdminData struct {
	AllUsers []dbmanagement.User
}

func Admin(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	adminData := AdminData{}
	sessionId, err := Session(w, r)
	if err != nil {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		fmt.Println("please log in as a user with admin permissions")
	}

	loggedInAs := dbmanagement.SelectUserFromSession(sessionId)
	if loggedInAs.Permission != "admin" {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		fmt.Println("please log in as a user with admin permissions")
	}

	if r.Method == "POST" {
		userToChange := r.FormValue("set to user")
		if userToChange != "" {
			dbmanagement.UpdateUserPermissionFromUUID(userToChange, "user")
		}
		userToChange = r.FormValue("set to moderator")
		if userToChange != "" {
			dbmanagement.UpdateUserPermissionFromUUID(userToChange, "moderator")
		}
		userToChange = r.FormValue("set to admin")
		if userToChange != "" {
			dbmanagement.UpdateUserPermissionFromUUID(userToChange, "admin")
		}
		userToChange = r.FormValue("delete user")
		if userToChange != "" {
			dbmanagement.DeleteFromTableWithUUID("Users", userToChange)
		}

	}
	if loggedInAs.Permission == "admin" {
		adminData.AllUsers = dbmanagement.SelectAllUsers()
		tmpl.ExecuteTemplate(w, "admin.html", adminData)
	}
}
