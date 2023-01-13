package controller

import (
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
	sessionId, err := Session(w, r)
	if err == nil {
		user := dbmanagement.SelectUserFromSession(sessionId)
		data.Cookie = sessionId
		data.UserInfo = user
		if r.Method == "POST" {
			comment := r.FormValue("post")
			if comment != "" {
				dbmanagement.InsertPost(comment, dbmanagement.SelectUserFromUUID(user.UUID).Name, 0, 0, "general", time.Now())
			}
		}
	}
	utils.HandleError("cant get user", err)
	posts := dbmanagement.SelectAllPosts()
	data.ListOfData = append(data.ListOfData, posts...)
	data.TitleName = "Forum"
	log.Println(data)
	tmpl.ExecuteTemplate(w, "forum.html", data)
}
