package controller

import (
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
	sessionId, err := Session(w, r)
	utils.HandleError("cant get sessionId:", err)
	posts := dbmanagement.SelectAllPosts()

	data := Data{}
	data.Cookie = sessionId
	data.ListOfData = append(data.ListOfData, posts...)
	tmpl.ExecuteTemplate(w, "forum.html", data)
}

func UsersPosts(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	sessionId, err := Session(w, r)
	utils.HandleError("cant get user", err)
	user := dbmanagement.SelectUserFromSession(sessionId)

	if r.Method == "POST" {
		comment := r.FormValue("post")
		if comment != "" {
			dbmanagement.InsertPost(comment, dbmanagement.SelectUserFromUUID(user.UUID).Name, 0, 0, "general", time.Now())
		}
	}
	posts := dbmanagement.SelectAllPosts()

	data := Data{}
	data.Cookie = sessionId
	data.UserInfo = user
	data.ListOfData = append(data.ListOfData, posts...)
	tmpl.ExecuteTemplate(w, "forum.html", data)
}
