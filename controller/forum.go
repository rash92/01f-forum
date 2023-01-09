package controller

import (
	"forum/dbmanagement"
	utils "forum/helpers"
	"html/template"
	"net/http"
	"time"
)

type Data struct {
	ListOfData []string
	Cookie     string
	UserInfo   string
}

// func UserLoggedIn(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
// 	if r.Method == "POST" {
// 		comment := r.FormValue("post")
// 		if comment != "" {
// 			dbmanagement.InsertPost(comment, "1", 0, 0, "general", time.Now())

// 		}
// 	}
// 	posts := dbmanagement.SelectAllPostsFromUser("1")
// 	data := Data{}
// 	for _, v := range posts {
// 		data.ListOfData = append(data.ListOfData, v.content)
// 	}
// 	tmpl.ExecuteTemplate(w, "user.html", data)
// }

func AllPosts(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {

	if r.Method == "POST" {
		comment := r.FormValue("post")
		if comment != "" {
			dbmanagement.InsertPost(comment, "1", 0, 0, "general", time.Now())

		}
	}
	posts := dbmanagement.SelectAllPosts()
	sessionId, err := Session(w, r)
	// log.Println("sessonID:", sessionId)
	utils.HandleError("cant get user", err)
	user := dbmanagement.SelectUserFromSession(sessionId)
	// log.Println("user is:", user)
	data := Data{}
	data.Cookie = sessionId
	data.UserInfo = user.Name
	for _, v := range posts {
		data.ListOfData = append(data.ListOfData, v.Content)
	}
	tmpl.ExecuteTemplate(w, "forum.html", data)
}
