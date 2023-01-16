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

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}

func GithubLogin(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	githubConfig := config.GithubSetupConfig()
	url := githubConfig.AuthCodeURL("randomstate")
	// redirect to github login page
	http.Redirect(w, r, url, http.StatusFound)
}

func GithubCallback(w http.ResponseWriter, r *http.Request, tmpl *template.Template) OauthAccount {
	//state
	state := r.URL.Query()["state"][0]
	if state != "randomstate" {
		fmt.Fprintln(w, "Google auth state error")
		return OauthAccount{}
	}

	// code
	// code := r.URL.Query()["code"][0]
	code := r.FormValue("code")
	fmt.Println("code is:", code)

	// configuration
	githubConfig := config.GithubSetupConfig()

	// exchange code for token
	token, err := githubConfig.Exchange(context.Background(), code)
	utils.HandleError("Code-taken exchange failed", err)

	client := githubConfig.Client(context.Background(), token)

	res, err := client.Get("https://api.github.com/user")
	utils.HandleError("Failed to fetch user data from github:", err)

	defer res.Body.Close()

	// parse response
	var value map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&value)
	utils.HandleError("Json parsing failed", err)

	account := OauthAccount{
		Name:  utils.AssertString(value["name"]),
		Email: utils.AssertString(value["email"]),
	}

	return account
}

func LoginUserWithGithub(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	account := GithubCallback(w, r, tmpl)

	user, err := dbmanagement.SelectUserFromEmail(account.Email)
	if err == nil {
		// create session cookie for user
		CreateUserSessionCookie(w, r, user)
	} else {
		// create user
		user := dbmanagement.InsertUser(account.Name, account.Email, "", "user")
		// create session cookie for user
		CreateUserSessionCookie(w, r, user)
		// utils.HandleError("Failed to create session in google authenticate", err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
