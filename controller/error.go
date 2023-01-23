package controller

import (
	"html/template"
	"net/http"
)

type Errors struct {
	Error string
}

func PageErrors(w http.ResponseWriter, r *http.Request, tmpl *template.Template, errortype string) {
	errors := Errors{
		Error: errortype,
	}
	tmpl.ExecuteTemplate(w, "error.html", errors)
}
