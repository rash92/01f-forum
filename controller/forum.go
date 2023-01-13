package controller

import (
	"fmt"
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Data struct {
	ListOfData []dbmanagement.Post
	Cookie     string
	UserInfo   dbmanagement.User
	TitleName  string
}

/*
Executes the forum.html template that includes all posts in the database.  SessionID is used the determine which user is currently using the website.

Also handles inserting a new post that updates in realtime.
*/
func AllPosts(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	data := Data{}
	sessionId, err := GetSessionIDFromBrowser(w, r)
	fmt.Println("session error is: ", err)
	user := dbmanagement.User{}
	if err == nil {
		user, err = dbmanagement.SelectUserFromSession(sessionId)
		data.Cookie = sessionId

		data.UserInfo = user
		fmt.Println("session id is: ", sessionId, "user info is: ", data.UserInfo, "cookie data is: ", data.Cookie)

		if r.Method == "POST" {
			comment := r.FormValue("post")
			tag := r.FormValue("tag")
			like := r.FormValue("like")
			dislike := r.FormValue("dislike")
			postid := r.FormValue("postid")
			if comment != "" {
				userFromUUID, err := dbmanagement.SelectUserFromUUID(user.UUID)
				utils.HandleError("cant get user with uuid in all posts", err)
				dbmanagement.InsertPost(comment, userFromUUID.Name, 0, 0, tag, time.Now())
				log.Println(tag)
				if !ExistingTag(tag) {
					dbmanagement.InsertTag(tag)
				}
			}
			if like == "Like" {
				dbmanagement.AddReactionToPost(user.UUID, postid, 1)
			}
			if dislike == "Dislike" {
				dbmanagement.AddReactionToPost(user.UUID, postid, -1)
			}
			idToDelete := r.FormValue("deletepost")
			fmt.Println("deleting post with id: ", idToDelete, " and contents: ", dbmanagement.SelectPostFromUUID(idToDelete))
			if idToDelete != "" {
				dbmanagement.DeleteFromTableWithUUID("Posts", idToDelete)
			}
		}

		utils.HandleError("cant get user", err)
		posts := dbmanagement.SelectAllPosts()
		for i, j := 0, len(posts)-1; i < j; i, j = i+1, j-1 {
			posts[i], posts[j] = posts[j], posts[i]
		}

		data := Data{}
		data.Cookie = sessionId
		data.UserInfo = user
		data.ListOfData = append(data.ListOfData, posts...)
		log.Println("Forum: ", data)
		tmpl.ExecuteTemplate(w, "forum.html", data)
		fmt.Println("user info is: ", data.UserInfo)
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
