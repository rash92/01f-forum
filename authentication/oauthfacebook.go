package auth

import (
	"context"
	"fmt"
	"forum/utils"
	"html/template"
	"net/http"
)

// Google Oauth
func FacebookLogin(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	facebookConfig := FacebookSetupConfig()
	url := facebookConfig.AuthCodeURL("randomstate")
	// redirect to facebook login page
	http.Redirect(w, r, url, http.StatusFound)
}

func FacebookCallback(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	// state
	state := r.FormValue("state")
	if state != "randomstate" {
		fmt.Fprintln(w, "facebook auth state error")
		return
	}

	// code
	code := r.FormValue("code")

	fmt.Println("code is:", code)
	// configuration
	facebookConfig := FacebookSetupConfig()

	// exchange code for token
	token, err := facebookConfig.Exchange(context.Background(), code)
	utils.HandleError("Code-taken exchange failed", err)

	// use google api to get user info
	resp, err := http.Get(facebookConfig.Endpoint.AuthURL + token.AccessToken)
	utils.HandleError("Failed to fetch user data from google:", err)

	fmt.Println("response:", resp)

	// parse response
	// value := ParseOauthResponse(resp)

	// fmt.Println("response:", value)

	// account := OauthAccount{
	// 	Name:  utils.AssertString(value["given_name"]),
	// 	Email: utils.AssertString(value["email"]),
	// }
	// login and create session for user
	// LoginUserWithOauth(w, r, tmpl, account)
}
