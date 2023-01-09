package controller

import (
	"forum/dbmanagement"
	"html/template"
	"net/http"
	"time"
)

type Data struct {
	ListOfData []string
}

func UserLoggedIn(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	if r.Method == "POST" {
		comment := r.FormValue("post")
		if comment != "" {
			dbmanagement.InsertPost(comment, 1, time.Now().String(), "general")
		}
	}
	posts := dbmanagement.DisplayAllPosts()
	data := Data{}
	for _, v := range posts {
		data.ListOfData = append(data.ListOfData, v.Comment)
	}
	tmpl.ExecuteTemplate(w, "user.html", data)
}
