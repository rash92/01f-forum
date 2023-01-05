package controller

import (
	"html/template"
	"net/http"
)

func Forum(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	tmpl.ExecuteTemplate(w, "forum.html", nil)
}
