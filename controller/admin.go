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
	AllUsers         []dbmanagement.User
	AdminRequests    []dbmanagement.AdminRequest
	ReportedPosts    []dbmanagement.Post
	ReportedComments []dbmanagement.Comment
	ReportedUsers    []dbmanagement.User
	TitleName        string
	UserInfo         dbmanagement.User
	TagsList         []dbmanagement.Tag
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
	utils.HandleError("cant get user", err)
	adminData.AllUsers = dbmanagement.SelectAllUsers()
	adminData.AdminRequests = dbmanagement.SelectAllAdminRequests()
	for _, adminRequest := range adminData.AdminRequests {
		if adminRequest.ReportedPostId != "" {
			adminData.ReportedPosts = append(adminData.ReportedPosts, dbmanagement.SelectPostFromUUID(adminRequest.ReportedPostId))
		}
		if adminRequest.ReportedCommentId != "" {
			adminData.ReportedComments = append(adminData.ReportedComments, dbmanagement.SelectCommentFromUUID(adminRequest.ReportedCommentId))
		}
		if adminRequest.ReportedUserId != "" {
			currentUser, err := dbmanagement.SelectUserFromUUID(adminRequest.ReportedUserId)
			utils.HandleError("couldn't select user when looking for reported users", err)
			adminData.ReportedUsers = append(adminData.ReportedUsers, currentUser)
		}
	}
	adminData.TitleName = "Admin"
	adminData.TagsList = dbmanagement.SelectAllTags()
	adminData.UserInfo = loggedInAs
	tmpl.ExecuteTemplate(w, "admin.html", adminData)
}
