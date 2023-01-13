package controller

import (
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"log"
	"net/http"
	"time"
)

type SubData struct {
	SubName    string
	ListOfData []dbmanagement.Post
	Cookie     string
	UserInfo   dbmanagement.User
}

func SubForum(w http.ResponseWriter, r *http.Request, tmpl *template.Template, tag string) {
	sessionId, err := Session(w, r)
	utils.HandleError("cant get user", err)
	user := dbmanagement.SelectUserFromSession(sessionId)

	if r.Method == "POST" {
		comment := r.FormValue("post")
		like := r.FormValue("like")
		dislike := r.FormValue("dislike")
		postid := r.FormValue("postid")
		if comment != "" {
			dbmanagement.InsertPost(comment, dbmanagement.SelectUserFromUUID(user.UUID).Name, 0, 0, tag, time.Now())
		}
		if like == "Like" {
			dbmanagement.AddReactionToPost(user.UUID, postid, 1)
		}
		if dislike == "Dislike" {
			dbmanagement.AddReactionToPost(user.UUID, postid, -1)
		}
	}
	posts := dbmanagement.SelectAllPostsFromTag(tag)
	for i, j := 0, len(posts)-1; i < j; i, j = i+1, j-1 {
		posts[i], posts[j] = posts[j], posts[i]
	}

	data := SubData{}
	data.SubName = "/" + tag
	data.Cookie = sessionId
	data.UserInfo = user
	data.ListOfData = append(data.ListOfData, posts...)
	log.Println("SubForum: ", data)
	tmpl.ExecuteTemplate(w, "subforum.html", data)
}
