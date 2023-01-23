package controller

import (
	"fmt"
	auth "forum/authentication"
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"net/http"
	"time"
)

type PostData struct {
	Post             dbmanagement.Post
	Comments         []dbmanagement.Comment
	NumOfComments    int
	Cookie           string
	UserInfo         dbmanagement.User
	TitleName        string
	HasNotifications bool
	Notifications    []dbmanagement.Notification
}

func Post(w http.ResponseWriter, r *http.Request, tmpl *template.Template, postid string) {
	data := Data{}
	sessionId, err := auth.GetSessionFromBrowser(w, r)
	fmt.Println("session error is: ", err)
	user := dbmanagement.User{}
	if err == nil {
		user, err = dbmanagement.SelectUserFromSession(sessionId)
		data.Cookie = sessionId

		data.UserInfo = user
		fmt.Println("session id is: ", sessionId, "user info is: ", data.UserInfo, "cookie data is: ", data.Cookie)

		if r.Method == "POST" {
			usertoken := dbmanagement.GetUserToken(user.UUID)
			if usertoken <= 0 {
				tmpl.ExecuteTemplate(w, "error.html ", nil)
			}
			dbmanagement.UpdateUserToken(user.UUID, 1)
			notfication := r.FormValue("notification")
			comment := r.FormValue("comment")
			like := r.FormValue("like")
			dislike := r.FormValue("dislike")
			commentlike := r.FormValue("commentlike")
			commentdislike := r.FormValue("commentdislike")
			if notfication != "" {
				dbmanagement.DeleteFromTableWithUUID("Notifications", notfication)
			}
			if comment != "" {
				userFromUUID, err := dbmanagement.SelectUserFromUUID(user.UUID)
				utils.HandleError("cant get user with uuid in all posts", err)
				thisComment := dbmanagement.InsertComment(comment, postid, userFromUUID.UUID, 0, 0, time.Now())
				post := dbmanagement.SelectPostFromUUID(postid)
				receiverId, _ := dbmanagement.SelectUserFromName(post.OwnerId)
				dbmanagement.AddNotification(receiverId.UUID, postid, thisComment.UUID, user.UUID, 0)
			}
			if like != "" {
				dbmanagement.AddReactionToPost(user.UUID, like, 1)
				post := dbmanagement.SelectPostFromUUID(like)
				receiverId, _ := dbmanagement.SelectUserFromName(post.OwnerId)
				dbmanagement.AddNotification(receiverId.UUID, like, "", user.UUID, 1)
			}
			if dislike != "" {
				dbmanagement.AddReactionToPost(user.UUID, dislike, -1)
				post := dbmanagement.SelectPostFromUUID(dislike)
				receiverId, _ := dbmanagement.SelectUserFromName(post.OwnerId)
				dbmanagement.AddNotification(receiverId.UUID, dislike, "", user.UUID, -1)
			}
			if commentlike != "" {
				dbmanagement.AddReactionToComment(user.UUID, commentlike, 1)
				comment := dbmanagement.SelectCommentFromUUID(commentlike)
				receiverId, _ := dbmanagement.SelectUserFromName(comment.OwnerId)
				dbmanagement.AddNotification(receiverId.UUID, "", commentlike, user.UUID, 1)
			}
			if commentdislike != "" {
				dbmanagement.AddReactionToComment(user.UUID, commentdislike, -1)
				comment := dbmanagement.SelectCommentFromUUID(commentdislike)
				receiverId, _ := dbmanagement.SelectUserFromName(comment.OwnerId)
				dbmanagement.AddNotification(receiverId.UUID, "", commentdislike, user.UUID, -1)
			}
			idToDelete := r.FormValue("deletepost")
			fmt.Println("deleting post with id: ", idToDelete, " and contents: ", dbmanagement.SelectPostFromUUID(idToDelete))
			if idToDelete != "" {
				dbmanagement.DeleteFromTableWithUUID("Posts", idToDelete)
			}
		}

		utils.HandleError("cant get user", err)
		post := dbmanagement.SelectPostFromUUID(postid)
		comments := dbmanagement.SelectAllCommentsFromPost(postid)
		for i, j := 0, len(comments)-1; i < j; i, j = i+1, j-1 {
			comments[i], comments[j] = comments[j], comments[i]
		}

		data := PostData{}
		data.Cookie = sessionId
		user.Notifications = dbmanagement.SelectAllNotificationsFromUser(user.UUID)
		data.UserInfo = user
		data.Post = post
		data.Comments = append(data.Comments, comments...)
		data.NumOfComments = len(comments)
		tmpl.ExecuteTemplate(w, "post.html", data)
	}
}
