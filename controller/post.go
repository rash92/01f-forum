package controller

import (
	"html/template"
	"net/http"
)

func Post(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	tmpl.ExecuteTemplate(w, "post.html", nil)
}
