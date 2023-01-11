package controller

import (
	"fmt"
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
}

/*
Executes the forum.html template that includes all posts in the database.  SessionID is used the determine which user is currently using the website.

Also handles inserting a new post that updates in realtime.
*/
func AllPosts(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	data := Data{}
	sessionId, err := Session(w, r)
	fmt.Println("session error is: ", err)
	user := dbmanagement.User{}
	if err == nil {
		user = dbmanagement.SelectUserFromSession(sessionId)
		data.Cookie = sessionId

		data.UserInfo = user
		fmt.Println("session id is: ", sessionId, "user info is: ", data.UserInfo, "cookie data is: ", data.Cookie)
		if r.Method == "POST" {
			comment := r.FormValue("post")
			fmt.Println(comment)
			if comment != "" {
				// post should be linked by the owners uuid and not the owners name i think?
				dbmanagement.InsertPost(comment, dbmanagement.SelectUserFromUUID(user.UUID).Name, 0, 0, "general", time.Now())
			}
			idToDelete := r.FormValue("deletepost")
			fmt.Println("deleting post with id: ", idToDelete, " and contents: ", dbmanagement.SelectPostFromUUID(idToDelete))
			if idToDelete != "" {
				dbmanagement.DeleteFromTableWithUUID("Posts", idToDelete)
			}

		}
	}
	utils.HandleError("cant get user", err)
	posts := dbmanagement.SelectAllPosts()
	data.ListOfData = append(data.ListOfData, posts...)
	tmpl.ExecuteTemplate(w, "forum.html", data)
	fmt.Println("user info is: ", data.UserInfo)
}
