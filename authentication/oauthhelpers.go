package auth

import (
	"context"
	"encoding/json"
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"net/http"

	"golang.org/x/oauth2"
)

func OauthToken(w http.ResponseWriter, r *http.Request, vendor string) (*oauth2.Token, error) {
	// state
	state := r.FormValue("state")
	if state != Randomstate {
		utils.WriteMessageToLogFile("Google auth state error")
	}

	// code
	code := r.FormValue("code")

	// configuration
	var oauthVendor *oauth2.Config

	if vendor == "google" {
		oauthVendor = GoogleSetupConfig()
	} else if vendor == "facebook" {
		oauthVendor = FacebookSetupConfig()
	} else if vendor == "github" {
		oauthVendor = GithubSetupConfig()
	}

	// exchange code for token
	token, err := oauthVendor.Exchange(context.Background(), code)
	utils.HandleError("Code-taken exchange failed", err)

	return token, err
}

// Helper to parse oauth response
func ParseOauthResponse(resp *http.Response) map[string]interface{} {
	var value map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&value)
	utils.HandleError("Json parsing failed", err)

	return value
}

// Gets account details from Oauth response and creates if they don't exist, if they do it creates a session and logs them in
func LoginUserWithOauth(w http.ResponseWriter, r *http.Request, tmpl *template.Template, account OauthAccount) {
	user, err := dbmanagement.SelectUserFromEmail(account.Email)
	_, err2 := dbmanagement.SelectUserFromName(account.Name)
	if err == nil && err2 == nil {
		// create session cookie for user
		CreateUserSession(w, r, user)
	} else {
		// create user
		user, _ := dbmanagement.InsertUser(account.Name, account.Email, "", "user", 0)
		// create session cookie for user
		CreateUserSession(w, r, user)
		// utils.HandleError("Failed to create session in google authenticate", err)
		LimitRequests(w, r, user)
		err = dbmanagement.UpdateUserToken(user.UUID, 1)
		if err != nil {
			http.Redirect(w, r, "/error", http.StatusSeeOther)
		}
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
