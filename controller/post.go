package controller

import (
	"forum/dbmanagement"
	"html/template"
	"net/http"
)

type PostData struct {
	Post      dbmanagement.Post
	Comments  []dbmanagement.Comment
	Cookie    string
	UserInfo  dbmanagement.User
	TitleName string
}

func Post(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	tmpl.ExecuteTemplate(w, "post.html", nil)
}
