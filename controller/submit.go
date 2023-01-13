package controller

import (
	"html/template"
	"net/http"
)

func Submit(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	data := Data{}
	data.TitleName = "Submit to Forum"
	tmpl.ExecuteTemplate(w, "submit.html", data)
}
