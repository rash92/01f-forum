package controller

import (
	"fmt"
	auth "forum/authentication"
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"net/http"
	"strings"
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
	TagsList         []dbmanagement.Tag
}

func CheckInputs(str string) bool {
	spl := strings.Fields(str)
	return len(spl) > 0
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
			notfication := r.FormValue("notification")
			comment := r.FormValue("comment")
			like := r.FormValue("like")
			dislike := r.FormValue("dislike")
			commentlike := r.FormValue("commentlike")
			commentdislike := r.FormValue("commentdislike")
			idToDelete := r.FormValue("deletepost")
			idToReport := r.FormValue("reportpost")
			deleteComment := r.FormValue("deletecomment")
			editComment := r.FormValue("editcomment")
			commentuuid := r.FormValue("commentuuid")

			if notfication != "" {
				dbmanagement.DeleteFromTableWithUUID("Notifications", notfication)
			}
			if CheckInputs(comment) {
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
			if idToDelete != "" {
				dbmanagement.DeletePostWithUUID(idToDelete)
			}
			if idToReport != "" {
				dbmanagement.CreateAdminRequest(user.UUID, user.Name, idToReport, "", "", "this post has been reported by a moderator")
			}
			if deleteComment != "" {
				dbmanagement.DeleteFromTableWithUUID("Comments", deleteComment)
			}
			if editComment != "" {
				dbmanagement.UpdateComment(commentuuid, editComment, postid, user.UUID, dbmanagement.SelectCommentFromUUID(editComment).Likes, dbmanagement.SelectCommentFromUUID(editComment).Dislikes, time.Now())
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
		data.TagsList = dbmanagement.SelectAllTags()
		tmpl.ExecuteTemplate(w, "post.html", data)
	}
}
