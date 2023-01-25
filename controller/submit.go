package controller

import (
	auth "forum/authentication"
	"forum/dbmanagement"
	"html/template"
	"log"
	"net/http"
)

func SubmitPost(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	data := Data{}
	user := dbmanagement.User{}
	sessionId, err := auth.GetSessionFromBrowser(w, r)
	log.Println("session error is: ", err)
	user, err = dbmanagement.SelectUserFromSession(sessionId)
	log.Println("session error is: ", err)
	data.TitleName = "Submit to Forum"
	data.Cookie = sessionId
	data.UserInfo = user
	data.TagsList = dbmanagement.SelectAllTags()
	tmpl.ExecuteTemplate(w, "submitpost.html", data)
}
