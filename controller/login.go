package controller

import (
	"html/template"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	tmpl.ExecuteTemplate(w, "index.html", nil)
}
