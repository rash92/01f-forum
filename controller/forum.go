package controller

import (
	"forum/dbmanagement"
	utils "forum/helpers"
	"html/template"
	"net/http"
	"time"
)

type Data struct {
	ListOfData []dbmanagement.Post
	Cookie     string
	UserInfo   dbmanagement.User
}

func AllPosts(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	sessionId, err := Session(w, r)
	// log.Println("sessonID:", sessionId)
	utils.HandleError("cant get user", err)
	user := dbmanagement.SelectUserFromSession(sessionId)

	if r.Method == "POST" {
		comment := r.FormValue("post")
		if comment != "" {
			dbmanagement.InsertPost(comment, dbmanagement.SelectUserFromUUID(user.UUID).Name, 0, 0, "general", time.Now())

		}
	}
	posts := dbmanagement.SelectAllPosts()

	// log.Println("user is:", user)
	data := Data{}
	data.Cookie = sessionId
	data.UserInfo = user
	data.ListOfData = append(data.ListOfData, posts...)
	tmpl.ExecuteTemplate(w, "forum.html", data)
}
