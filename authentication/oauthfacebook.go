package auth

import (
	"context"
	"forum/utils"
	"html/template"
	"net/http"
)

// Google Oauth
func FacebookLogin(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	facebookConfig := FacebookSetupConfig()
	url := facebookConfig.AuthCodeURL(Randomstate)
	// redirect to facebook login page
	http.Redirect(w, r, url, http.StatusFound)
}

func FacebookCallback(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	// state
	state := r.FormValue("state")
	if state != Randomstate {
		utils.WriteMessageToLogFile("Fcaebook auth state error")
		return
	}

	// code
	code := r.FormValue("code")

	// configuration
	facebookConfig := FacebookSetupConfig()

	// exchange code for token
	token, err := facebookConfig.Exchange(context.Background(), code)
	utils.HandleError("Code-taken exchange failed", err)

	// use google api to get user info
	resp, err := http.Get(FBAuthURL + token.AccessToken)
	utils.HandleError("Failed to fetch user data from facebook:", err)

	defer resp.Body.Close()
	// parse response
	value := ParseOauthResponse(resp)

	account := OauthAccount{
		Name:  utils.AssertString(value["name"]),
		Email: utils.AssertString(value["email"]),
	}
	// login and create session for user
	LoginUserWithOauth(w, r, tmpl, account)
}
