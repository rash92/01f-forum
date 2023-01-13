package controller

import (
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"net/http"
	"time"
)

type Data struct {
	ListOfData []dbmanagement.Post
	Cookie     string
	UserInfo   dbmanagement.User
}

/*
Executes the forum.html template that includes all posts in the database.  SessionID is used the determine which user is currently using the website.

Also handles inserting a new post that updates in realtime.
*/
// func AllPosts(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
// 	sessionId, err := GetSessionIDFromBrowser(w, r)
// 	utils.HandleError("Cannot get sessionId:", err)
// 	posts := dbmanagement.SelectAllPosts()

// 	data := Data{}
// 	data.Cookie = sessionId
// 	data.ListOfData = append(data.ListOfData, posts...)
// 	tmpl.ExecuteTemplate(w, "forum.html", data)
// }

// func UsersPosts(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
// 	sessionId, err := GetSessionIDFromBrowser(w, r)
// 	utils.HandleError("Cannot get sessionID for users posts", err)

// 	user, err := dbmanagement.SelectUserFromSession(sessionId)
// 	utils.HandleError("Cannot get user with sessionID", err)

// 	if r.Method == "POST" {
// 		comment := r.FormValue("post")
// 		if comment != "" {
// 			userByUUID, err := dbmanagement.SelectUserFromUUID(user.UUID)
// 			utils.HandleError("Cannot get user by UUID for comment:", err)

// 			dbmanagement.InsertPost(comment, userByUUID.Name, 0, 0, "general", time.Now())
// 		}
// 	}
// 	posts := dbmanagement.SelectAllPosts()

// 	data := Data{}
// 	data.Cookie = sessionId
// 	data.UserInfo = user
// 	data.ListOfData = append(data.ListOfData, posts...)

//		tmpl.ExecuteTemplate(w, "forum.html", data)
//	}
func AllPosts(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	data := Data{}
	sessionId, err := GetSessionIDFromBrowser(w, r)

	if err == nil {
		user, err := dbmanagement.SelectUserFromSession(sessionId)
		utils.HandleError("Cannot get user with sessionID", err)
		data.Cookie = sessionId
		data.UserInfo = user
		if r.Method == "POST" {
			comment := r.FormValue("post")
			if comment != "" {
				userByUUID, err := dbmanagement.SelectUserFromUUID(user.UUID)
				utils.HandleError("Cannot get user by UUID for comment:", err)

				dbmanagement.InsertPost(comment, userByUUID.Name, 0, 0, "general", time.Now())
			}
		}
	}

	posts := dbmanagement.SelectAllPosts()
	data.ListOfData = append(data.ListOfData, posts...)
	tmpl.ExecuteTemplate(w, "forum.html", data)
}
