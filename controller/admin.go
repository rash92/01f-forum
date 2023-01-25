package controller

import (
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
}

// username: admin password: admin for existing user with admin permissions, can create and change other users to be admin while logged in as anyone who is admin
func Admin(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	adminData := AdminData{}
	sessionId, err := auth.GetSessionFromBrowser(w, r)
	utils.HandleError("Unable to get session from browser in admin handler", err)
	user, err := dbmanagement.SelectUserFromSession(sessionId)
	utils.HandleError("Unable to user session from admin handler", err)
	if err != nil {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		utils.WriteMessageToLogFile("Please log in as a user with admin permissions")
		return
	}

	loggedInAs, err := dbmanagement.SelectUserFromSession(sessionId)
	utils.HandleError("Unable get logged in user in admin", err)
	if loggedInAs.Permission != "admin" {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		utils.WriteMessageToLogFile("Please log in as a user with admin permissions")
		return
	}

	if r.Method == "POST" {
		err := dbmanagement.UpdateUserToken(user.UUID, 1)
		if err != nil {
			http.Redirect(w, r, "/error", http.StatusSeeOther)
			return
		}
		dbmanagement.UpdateUserToken(user.UUID, 1)
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
