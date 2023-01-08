package controller

import (
	"html/template"
	"net/http"
)

func UserLoggedIn(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	tmpl.ExecuteTemplate(w, "user.html", nil)
}
