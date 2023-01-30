package controller

import (
	"fmt"
	auth "forum/authentication"
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"net/http"
)

type UserData struct {
	UserPosts         []dbmanagement.Post
	LikedUserPosts    []dbmanagement.Post
	DislikedUserPosts []dbmanagement.Post
	UserComments      []dbmanagement.Comment
	UserInfo          dbmanagement.User
	TitleName         string
	Cookie            string
	TagsList          []dbmanagement.Tag
}

func User(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	auth.LoggedInStatus(w, r, tmpl)
	data := UserData{}
	SessionId, err := auth.GetSessionFromBrowser(w, r)
	data.UserInfo, err = dbmanagement.SelectUserFromSession(SessionId)
	if err != nil {
		utils.HandleError("couldn't find user sessions id", err)
	}

	if SessionId == "" {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		// fmt.Println("please log in")
		return
	}

	if r.Method == "POST" {
		postIdToDelete := r.FormValue("deletepost")
		// fmt.Println("deleting post with id: ", postIdToDelete, " and contents: ", dbmanagement.SelectPostFromUUID(postIdToDelete))
		if postIdToDelete != "" {
			dbmanagement.DeletePostWithUUID(postIdToDelete)
		}
		commentIdToDelete := r.FormValue("deletecomment")
		// fmt.Println("deleting comment with id: ", commentIdToDelete, " and contents: ", dbmanagement.SelectCommentFromUUID(commentIdToDelete))
		if commentIdToDelete != "" {
			dbmanagement.DeleteFromTableWithUUID("Comments", commentIdToDelete)
		}
		userIdToRequestModerator := r.FormValue("request to become moderator")
		// fmt.Println("requesting user id: ", userIdToRequestModerator, "to become moderator")
		if userIdToRequestModerator != "" {
			newrequest := dbmanagement.CreateAdminRequest(userIdToRequestModerator, data.UserInfo.Name, "", "", "", "this user is asking to become a moderator")
			fmt.Println("new request content is: ", newrequest.Description)
		}

	}

	data.UserInfo.Notifications = dbmanagement.SelectAllNotificationsFromUser(data.UserInfo.UUID)
	utils.HandleError("Could not get user session in user", err)
	data.UserPosts = dbmanagement.SelectAllPostsFromUser(data.UserInfo.Name)
	data.LikedUserPosts = dbmanagement.SelectAllLikedPostsFromUser(data.UserInfo)
	data.DislikedUserPosts = dbmanagement.SelectAllDislikedPostsFromUser(data.UserInfo)
	data.UserComments = dbmanagement.SelectAllCommentsFromUser(data.UserInfo.UUID)
	data.TitleName = data.UserInfo.Name
	data.TagsList = dbmanagement.SelectAllTags()
	tmpl.ExecuteTemplate(w, "user.html", data)
}
