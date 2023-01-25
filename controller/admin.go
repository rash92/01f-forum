package controller

import (
	"fmt"
	auth "forum/authentication"
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"net/http"
)

type AdminData struct {
	AllUsers      []dbmanagement.User
	AllTags       []dbmanagement.Tag
	AdminRequests []dbmanagement.AdminRequest
	TitleName     string
	UserInfo      dbmanagement.User
}

// username: admin password: admin for existing user with admin permissions, can create and change other users to be admin while logged in as anyone who is admin
func Admin(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	adminData := AdminData{}
	sessionId, err := auth.GetSessionFromBrowser(w, r)
	if err != nil {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		fmt.Println("please log in as a user with admin permissions")
		return
	}

	loggedInAs, err := dbmanagement.SelectUserFromSession(sessionId)
	utils.HandleError("cant get logged in user in admin", err)
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
		tagToDeletePostsLinkedTo := r.FormValue("delete all posts with tag")
		if tagToDeletePostsLinkedTo != "" {
			dbmanagement.DeleteAllPostsWithTag(tagToDeletePostsLinkedTo)
		}
		adminRequestToDelete := r.FormValue("delete request")
		if adminRequestToDelete != "" {
			dbmanagement.DeleteFromTableWithUUID("AdminRequests", adminRequestToDelete)
		}

	}
	user, err := dbmanagement.SelectUserFromSession(sessionId)
	utils.HandleError("cant get user", err)
	adminData.AllUsers = dbmanagement.SelectAllUsers()
	adminData.AdminRequests = dbmanagement.SelectAllAdminRequests()
	adminData.AllTags = dbmanagement.SelectAllTags()
	adminData.TitleName = "Admin"
	adminData.UserInfo = user
	tmpl.ExecuteTemplate(w, "admin.html", adminData)
}
