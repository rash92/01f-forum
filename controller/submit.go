package controller

import (
	auth "forum/authentication"
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"net/http"
)

func SubmitPost(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	data := Data{}
	user := dbmanagement.User{}
	sessionId, err := auth.GetSessionFromBrowser(w, r)
	utils.HandleError("Unable to get session from browser in SubmitPost function", err)
	user, err = dbmanagement.SelectUserFromSession(sessionId)
	utils.HandleError("Unable to select user with sessionID in SubmitPost function", err)
	data.TitleName = "Submit to Forum"
	data.Cookie = sessionId
	data.UserInfo = user
	tmpl.ExecuteTemplate(w, "submitpost.html", data)
}
