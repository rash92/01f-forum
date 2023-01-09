package controller

import (
	"fmt"
	"forum/dbmanagement"
	"html/template"
	"net/http"
	"time"
)

type Data struct {
	ListOfData []string
}

func Forum(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	user := r.FormValue("user_name")
	fmt.Println("This is it: " + user)
	if r.Method == "POST" {
		comment := r.FormValue("post")
		if comment != "" {
			dbmanagement.InsertPost(comment, 1, time.Now().String(), "general")
		}
	}
	posts := dbmanagement.DisplayAllPosts()
	data := Data{}
	for _, v := range posts {
		data.ListOfData = append(data.ListOfData, v.PostText)
	}
	tmpl.ExecuteTemplate(w, "forum.html", data)
}
