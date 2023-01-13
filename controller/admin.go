package controller

import (
	"fmt"
	"forum/dbmanagement"
	"html/template"
	"net/http"
)

type AdminData struct {
	AllUsers      []dbmanagement.User
	AllTags       []dbmanagement.Tag
	AdminRequests []dbmanagement.AdminRequest
}

// username: admin password: admin for existing user with admin permissions, can create and change other users to be admin while logged in as anyone who is admin
func Admin(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	adminData := AdminData{}
	sessionId, err := Session(w, r)
	if err != nil {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		fmt.Println("please log in as a user with admin permissions")
		return
	}

	loggedInAs := dbmanagement.SelectUserFromSession(sessionId)
	if loggedInAs.Permission != "admin" {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		fmt.Println("please log in as a user with admin permissions")
		return
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
		tagToChange := r.FormValue("delete tag")
		if tagToChange != "" {
			dbmanagement.DeleteFromTableWithUUID("Tags", tagToChange)
		}
		adminRequestToDelete := r.FormValue("delete request")
		if adminRequestToDelete != "" {
			dbmanagement.DeleteFromTableWithUUID("AdminRequests", adminRequestToDelete)
		}

	}

	adminData.AllUsers = dbmanagement.SelectAllUsers()
	adminData.AdminRequests = dbmanagement.SelectAllAdminRequests()
	adminData.AllTags = dbmanagement.SelectAllTags()
	tmpl.ExecuteTemplate(w, "admin.html", adminData)
}
