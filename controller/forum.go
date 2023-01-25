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

type Data struct {
	ListOfData []dbmanagement.Post
	Cookie     string
	UserInfo   dbmanagement.User
	TitleName  string
	IsCorrect  bool
}

/*
Executes the forum.html template that includes all posts in the database.  SessionID is used the determine which user is currently using the website.

Also handles inserting a new post that updates in realtime.
*/
func AllPosts(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	data := Data{}
	sessionId, err := auth.GetSessionFromBrowser(w, r)
	// fmt.Println("session error is: ", err)
	if sessionId == "" {
		err := auth.CreateUserSession(w, r, dbmanagement.User{})
		if err != nil {
			utils.HandleError("Unable to create visitor session", err)
		} else {
			sessionId, _ = auth.GetSessionFromBrowser(w, r)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}

	user := dbmanagement.User{}
	if err == nil {
		user, err = dbmanagement.SelectUserFromSession(sessionId)
		data.Cookie = sessionId
		filterOrder := false
		data.UserInfo = user
		fmt.Println("session id is: ", sessionId, "user info is: ", data.UserInfo, "cookie data is: ", data.Cookie)

		if r.Method == "POST" {
			err := dbmanagement.UpdateUserToken(user.UUID, 1)
			if err != nil {
				http.Redirect(w, r, "/error", http.StatusSeeOther)
			}

			title := r.FormValue("submission-title")
			content := r.FormValue("post")
			tag := r.FormValue("tag")
			like := r.FormValue("like")
			dislike := r.FormValue("dislike")
			filter := r.FormValue("filter")
			if filter == "oldest" {
				filterOrder = true
			}
			if content != "" {
				userFromUUID, err := dbmanagement.SelectUserFromUUID(user.UUID)
				utils.HandleError("Unable to get user with uuid", err)
				dbmanagement.InsertPost(title, content, userFromUUID.Name, 0, 0, tag, time.Now())
				// log.Println(tag)
				if !ExistingTag(tag) {
					dbmanagement.InsertTag(tag)
				}
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

			idToDelete := r.FormValue("deletepost")
			// fmt.Println("deleting post with id: ", idToDelete, " and contents: ", dbmanagement.SelectPostFromUUID(idToDelete))
			if idToDelete != "" {
				dbmanagement.DeleteFromTableWithUUID("Posts", idToDelete)
			}
		}

		utils.HandleError("Unable to get user", err)
		posts := dbmanagement.SelectAllPosts()
		if !filterOrder {
			for i, j := 0, len(posts)-1; i < j; i, j = i+1, j-1 {
				posts[i], posts[j] = posts[j], posts[i]
			}
		}

		data := Data{}
		data.Cookie = sessionId
		user.Notifications = dbmanagement.SelectAllNotificationsFromUser(user.UUID)
		data.UserInfo = user
		// log.Println("SESSION ID: ", data.Cookie)
		// log.Println("CURRENT USER: ", data.UserInfo.Name)
		data.ListOfData = append(data.ListOfData, posts...)
		// fmt.Println("Forum data: ", data)
		tmpl.ExecuteTemplate(w, "forum.html", data)
	}
}

func ExistingTag(tag string) bool {
	allTags := dbmanagement.SelectAllTags()
	for _, v := range allTags {
		if tag == v.TagName {
			return true
		}
	}
	return false
}
