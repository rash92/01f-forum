package controller

import (
	"context"
	"fmt"
	"forum/config"
	"forum/utils"
	"html/template"
	"io"
	"net/http"
)

func GoogleLogin(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	googleConfig := config.SetupConfig()
	url := googleConfig.AuthCodeURL("randomstate")
	//redirect to google login page
	http.Redirect(w, r, url, http.StatusFound)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	// state
	state := r.URL.Query()["state"][0]
	if state != "randomstate" {
		fmt.Fprintln(w, "Google auth state error")
		return
	}

	// code
	code := r.URL.Query()["code"][0]

	// configuration
	googleConfig := config.SetupConfig()

	//exchange code for token
	token, err := googleConfig.Exchange(context.Background(), code)
	utils.HandleError("Code-taken exchange failed", err)

	// use google api to get user info
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	utils.HandleError("Failed to fetch user data from google:", err)

	// parse response
	userData, err := io.Copy(w, resp.Body)
	utils.HandleError("Json parsing failed", err)

	fmt.Println(w, userData)

}
