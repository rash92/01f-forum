package controller

import (
	"html/template"
	"net/http"
)

/*
Displays the register page...
*/
func Register(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	data := Data{}
	data.TitleName = "Sign Up"
	tmpl.ExecuteTemplate(w, "register.html", data)
}
