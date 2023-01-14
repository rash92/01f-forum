package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"forum/config"
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"net/http"
)

type GoogleAccount struct {
	Name, Email string
}

func GoogleLogin(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	googleConfig := config.SetupConfig()
	url := googleConfig.AuthCodeURL("randomstate")
	// redirect to google login page
	http.Redirect(w, r, url, http.StatusFound)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request, tmpl *template.Template) GoogleAccount {
	// state
	state := r.URL.Query()["state"][0]
	if state != "randomstate" {
		fmt.Fprintln(w, "Google auth state error")
		return GoogleAccount{}
	}

	// code
	code := r.URL.Query()["code"][0]

	// configuration
	googleConfig := config.SetupConfig()

	// exchange code for token
	token, err := googleConfig.Exchange(context.Background(), code)
	utils.HandleError("Code-taken exchange failed", err)

	// use google api to get user info
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	utils.HandleError("Failed to fetch user data from google:", err)

	// parse response
	var value map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&value)
	utils.HandleError("Json parsing failed", err)

	account := GoogleAccount{
		Name:  utils.AssertString(value["given_name"]),
		Email: utils.AssertString(value["email"]),
	}

	return account
}

func LoginUserWithGoogle(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	account := GoogleCallback(w, r, tmpl)

	user, err := dbmanagement.SelectUserFromEmail(account.Email)
	if err == nil {
		// create session cookie for user
		CreateUserSessionCookie(w, r, user)
	} else {
		// create user
		dbmanagement.InsertUser(account.Name, account.Email, "", "user")
		// create session cookie for user
		CreateUserSessionCookie(w, r, user)
		// utils.HandleError("Failed to create session in google authenticate", err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
