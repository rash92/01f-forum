package controller

import (
	"html/template"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	tmpl.ExecuteTemplate(w, "register.html", nil)
}
