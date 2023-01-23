package auth

import (
	"encoding/json"
	"forum/dbmanagement"
	"forum/utils"
	"html/template"
	"net/http"
)

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
	if err == nil {
		// create session cookie for user
		CreateUserSession(w, r, user)
	} else {
		// create user
		user, _ := dbmanagement.InsertUser(account.Name, account.Email, "", "user", 0)
		// create session cookie for user
		CreateUserSession(w, r, user)
		// utils.HandleError("Failed to create session in google authenticate", err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
