package controller

import (
	"fmt"
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"net/http"
)

type UserData struct {
	UserPosts    []dbmanagement.Post
	UserComments []dbmanagement.Comment
	User         dbmanagement.User
}

func User(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	data := UserData{}
	SessionId, err := GetSessionIDFromBrowser(w, r)
	if err != nil {
		utils.HandleError("couldn't find user sessions id", err)
	}

	if SessionId == "" {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		fmt.Println("please log in")
		return
	}

	data.User, err = dbmanagement.SelectUserFromSession(SessionId)
	data.UserPosts = dbmanagement.SelectAllPostsFromUser(data.User.Name)
	data.UserComments = dbmanagement.SelectAllCommentsFromUser(data.User.Name)

	if r.Method == "POST" {
		postIdToDelete := r.FormValue("deletepost")
		fmt.Println("deleting post with id: ", postIdToDelete, " and contents: ", dbmanagement.SelectPostFromUUID(postIdToDelete))
		if postIdToDelete != "" {
			dbmanagement.DeleteFromTableWithUUID("Posts", postIdToDelete)
		}
		commentIdToDelete := r.FormValue("deletecomment")
		fmt.Println("deleting comment with id: ", commentIdToDelete, " and contents: ", dbmanagement.SelectCommentFromUUID(commentIdToDelete))
		if postIdToDelete != "" {
			dbmanagement.DeleteFromTableWithUUID("Comments", commentIdToDelete)
		}
		userIdToRequestModerator := r.FormValue("request to become moderator")
		fmt.Println("requesting user id: ", userIdToRequestModerator, "to become moderator")
		if userIdToRequestModerator != "" {
			newrequest := dbmanagement.CreateAdminRequest(userIdToRequestModerator, data.User.Name, "this user is asking to become a moderator")
			fmt.Println("new request content is: ", newrequest.Content)
		}
	}

	tmpl.ExecuteTemplate(w, "user.html", data)
}
