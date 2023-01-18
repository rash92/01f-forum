package controller

import (
	"fmt"
	auth "forum/authentication"
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Data struct {
	ListOfData []dbmanagement.Post
	Cookie     string
	UserInfo   dbmanagement.User
	TitleName  string
}

/*
Executes the forum.html template that includes all posts in the database.  SessionID is used the determine which user is currently using the website.

Also handles inserting a new post that updates in realtime.
*/
func AllPosts(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	data := Data{}
	sessionId, err := auth.GetSessionFromBrowser(w, r)
	// fmt.Println("session error is: ", err)
	if sessionId == "" {
		err := auth.CreateUserSession(w, r, dbmanagement.User{})
		if err != nil {
			utils.HandleError("unable to create visitor session", err)
		} else {
			sessionId, _ = auth.GetSessionFromBrowser(w, r)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}

	user := dbmanagement.User{}
	if err == nil {
		user, err = dbmanagement.SelectUserFromSession(sessionId)
		data.Cookie = sessionId

		data.UserInfo = user
		// fmt.Println("session id is: ", sessionId, "user info is: ", data.UserInfo, "cookie data is: ", data.Cookie)

		if r.Method == "POST" {
			SubmissionHandler(w, r, user)
			// comment := r.FormValue("post")
			// tags := r.FormValue("tag")
			// like := r.FormValue("like")
			// dislike := r.FormValue("dislike")

			// if comment != "" {
			// 	userFromUUID, err := dbmanagement.SelectUserFromUUID(user.UUID)
			// 	utils.HandleError("cant get user with uuid in all posts", err)
			// 	post := dbmanagement.InsertPost(comment, userFromUUID.Name, 0, 0, time.Now())
			// 	// log.Println(tag)

			// 	if tags != "" {
			// 		tagslice := strings.Fields(tags)
			// 		for _, tagname := range tagslice {
			// 			if !ExistingTag(tagname) {
			// 				dbmanagement.InsertTag(tagname)
			// 			}
			// 			tag, err := dbmanagement.SelectTagFromName(tagname)
			// 			utils.HandleError("unable to retrieve tag id", err)
			// 			dbmanagement.InsertTaggedPost(tag.UUID, post.UUID)
			// 		}
			// 	}
			// }

			// if like != "" {
			// 	dbmanagement.AddReactionToPost(user.UUID, like, 1)
			// }
			// if dislike != "" {
			// 	dbmanagement.AddReactionToPost(user.UUID, dislike, -1)
			// }

			// idToDelete := r.FormValue("deletepost")
			// // fmt.Println("deleting post with id: ", idToDelete, " and contents: ", dbmanagement.SelectPostFromUUID(idToDelete))
			// if idToDelete != "" {
			// 	dbmanagement.DeleteFromTableWithUUID("Posts", idToDelete)
			// }
		}

		utils.HandleError("cant get user", err)
		posts := dbmanagement.SelectAllPosts()
		for i, j := 0, len(posts)-1; i < j; i, j = i+1, j-1 {
			posts[i], posts[j] = posts[j], posts[i]
		}

		data := Data{}
		data.Cookie = sessionId
		data.UserInfo = user
		data.ListOfData = append(data.ListOfData, posts...)
		// fmt.Println("Forum data: ", data)
		tmpl.ExecuteTemplate(w, "forum.html", data)
	}
}

func ExistingTag(tag string) bool {
	allTags := dbmanagement.SelectAllTags()
	for _, v := range allTags {
		if tag == v.TagName {
			return true
		}
	}
	return false
}

// followed this: https://freshman.tech/file-upload-golang/
func SubmissionHandler(w http.ResponseWriter, r *http.Request, user dbmanagement.User) {
	// 20 megabytes
	maxSize := 20 * 1024 * 1024

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxSize))
	err := r.ParseMultipartForm(int64(maxSize))
	if err != nil {
		utils.HandleError("error parsing form for image, likely too big", err)
		http.Error(w, "max filesize is 20Mb, please upload a smaller image", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("submission-image")
	if err != nil {
		utils.HandleError("error retrieving file from form", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		utils.HandleError("error creating file directory for uploads", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	destinationFile, err := os.Create(fmt.Sprint("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
	if err != nil {
		utils.HandleError("error creating file for image", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, file)

	if err != nil {
		utils.HandleError("error copying file to destination", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "file upload successful")

	comment := r.FormValue("post")
	tags := r.FormValue("tag")
	like := r.FormValue("like")
	dislike := r.FormValue("dislike")

	if comment != "" {
		userFromUUID, err := dbmanagement.SelectUserFromUUID(user.UUID)
		utils.HandleError("cant get user with uuid in all posts", err)
		post := dbmanagement.InsertPost(comment, userFromUUID.Name, 0, 0, time.Now())
		// log.Println(tag)

		if tags != "" {
			tagslice := strings.Fields(tags)
			for _, tagname := range tagslice {
				if !ExistingTag(tagname) {
					dbmanagement.InsertTag(tagname)
				}
				tag, err := dbmanagement.SelectTagFromName(tagname)
				utils.HandleError("unable to retrieve tag id", err)
				dbmanagement.InsertTaggedPost(tag.UUID, post.UUID)
			}
		}
	}

	if like != "" {
		dbmanagement.AddReactionToPost(user.UUID, like, 1)
	}
	if dislike != "" {
		dbmanagement.AddReactionToPost(user.UUID, dislike, -1)
	}

	idToDelete := r.FormValue("deletepost")
	// fmt.Println("deleting post with id: ", idToDelete, " and contents: ", dbmanagement.SelectPostFromUUID(idToDelete))
	if idToDelete != "" {
		dbmanagement.DeleteFromTableWithUUID("Posts", idToDelete)
	}
}
