package controller

import (
	"html/template"
	"net/http"
)

func Submit(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	tmpl.ExecuteTemplate(w, "submit.html", nil)
}
