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
	Post          dbmanagement.Post
	Comments      []dbmanagement.Comment
	NumOfComments int
	Cookie        string
	UserInfo      dbmanagement.User
	TitleName     string
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
			comment := r.FormValue("comment")
			like := r.FormValue("like")
			dislike := r.FormValue("dislike")
			commentlike := r.FormValue("commentlike")
			commentdislike := r.FormValue("commentdislike")
			if comment != "" {
				userFromUUID, err := dbmanagement.SelectUserFromUUID(user.UUID)
				utils.HandleError("cant get user with uuid in all posts", err)
				dbmanagement.InsertComment(comment, postid, userFromUUID.UUID, 0, 0, time.Now())
			}
			if like != "" {
				dbmanagement.AddReactionToPost(user.UUID, like, 1)
			}
			if dislike != "" {
				dbmanagement.AddReactionToPost(user.UUID, dislike, -1)
			}
			if commentlike != "" {
				dbmanagement.AddReactionToComment(user.UUID, commentlike, 1)
			}
			if commentdislike != "" {
				dbmanagement.AddReactionToComment(user.UUID, commentdislike, -1)
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
		data.UserInfo = user
		data.Post = post
		data.Comments = append(data.Comments, comments...)
		data.NumOfComments = len(comments)
		tmpl.ExecuteTemplate(w, "post.html", data)
	}
}
