package controller

import (
	"html/template"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	cookie := http.Cookie{
		Name:     "_cookie",
		Value:    "",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/forum", http.StatusFound)
}
