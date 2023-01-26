package controller

import (
	auth "forum/authentication"
	"forum/dbmanagement"
	"html/template"
	"log"
	"net/http"
)

type SubmitData struct {
	ListOfData []dbmanagement.Post
	Cookie     string
	UserInfo   dbmanagement.User
	TitleName  string
	IsCorrect  bool
	IsEdit     bool
	EditPost   dbmanagement.Post
	Tags       string
	TagsList   []dbmanagement.Tag
}

func SubmitPost(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	auth.LoggedInStatus(w, r, tmpl)
	data := SubmitData{}
	user := dbmanagement.User{}
	tags := []dbmanagement.Tag{}
	sessionId, err := auth.GetSessionFromBrowser(w, r)
	log.Println("session error is: ", err)
	user, err = dbmanagement.SelectUserFromSession(sessionId)
	log.Println("session error is: ", err)
	if r.Method == "POST" {
		idToEdit := r.FormValue("editpost")
		if idToEdit != "" {
			data.IsEdit = true
			data.EditPost = dbmanagement.SelectPostFromUUID(idToEdit)
			tags = dbmanagement.SelectAllTagsFromPost(data.EditPost.UUID)
		}
	}
	data.TitleName = "Submit to Forum"
	data.Cookie = sessionId
	data.UserInfo = user
	tagsAsString := ""
	for _, v := range tags {
		tagsAsString += v.TagName
		tagsAsString += " "
	}
	data.Tags = tagsAsString
	data.TagsList = dbmanagement.SelectAllTags()
	tmpl.ExecuteTemplate(w, "submitpost.html", data)
}
