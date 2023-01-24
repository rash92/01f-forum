package controller

import (
	"fmt"
	auth "forum/authentication"
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"log"
	"net/http"
)

type UserData struct {
	UserPosts    []dbmanagement.Post
	UserComments []dbmanagement.Comment
	UserInfo     dbmanagement.User
	TitleName    string
	Cookie       string
}

func User(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	data := UserData{}
	SessionId, err := auth.GetSessionFromBrowser(w, r)
	if err != nil {
		utils.HandleError("couldn't find user sessions id", err)
	}

	if SessionId == "" {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		// fmt.Println("please log in")
		return
	}

	data.UserInfo, err = dbmanagement.SelectUserFromSession(SessionId)
	data.UserInfo.Notifications = dbmanagement.SelectAllNotificationsFromUser(data.UserInfo.UUID)
	utils.HandleError("Could not get user session in user", err)
	data.UserPosts = dbmanagement.SelectAllPostsFromUser(data.UserInfo.Name)
	data.UserComments = dbmanagement.SelectAllCommentsFromUser(data.UserInfo.Name)
	log.Println(data.UserComments)
	data.TitleName = "Welcome"

	if r.Method == "POST" {
		err := dbmanagement.UpdateUserToken(data.UserInfo.UUID, 1)
		if err != nil {
			http.Redirect(w, r, "/rate_error", http.StatusTooManyRequests)
			return
		}
		postIdToDelete := r.FormValue("deletepost")
		// fmt.Println("deleting post with id: ", postIdToDelete, " and contents: ", dbmanagement.SelectPostFromUUID(postIdToDelete))
		if postIdToDelete != "" {
			dbmanagement.DeleteFromTableWithUUID("Posts", postIdToDelete)
		}
		commentIdToDelete := r.FormValue("deletecomment")
		// fmt.Println("deleting comment with id: ", commentIdToDelete, " and contents: ", dbmanagement.SelectCommentFromUUID(commentIdToDelete))
		if postIdToDelete != "" {
			dbmanagement.DeleteFromTableWithUUID("Comments", commentIdToDelete)
		}
		userIdToRequestModerator := r.FormValue("request to become moderator")
		// fmt.Println("requesting user id: ", userIdToRequestModerator, "to become moderator")
		if userIdToRequestModerator != "" {
			newrequest := dbmanagement.CreateAdminRequest(userIdToRequestModerator, data.UserInfo.Name, "this user is asking to become a moderator")
			fmt.Println("new request content is: ", newrequest.Content)
		}
	}

	tmpl.ExecuteTemplate(w, "user.html", data)
}
