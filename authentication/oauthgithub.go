package auth

import (
	"context"
	"fmt"
	"forum/utils"
	"html/template"
	"net/http"
)

// Github Oauth
func GithubLogin(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	githubConfig := GithubSetupConfig()
	url := githubConfig.AuthCodeURL(Randomstate)
	// redirect to github login page
	http.Redirect(w, r, url, http.StatusFound)
}

func GithubCallback(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	//state
	state := r.FormValue("state")
	if state != Randomstate {
		utils.WriteMessageToLogFile("Github auth state error")
		return
	}

	// code
	code := r.FormValue("code")
	// configuration
	githubConfig := GithubSetupConfig()

	// exchange code for token
	token, err := githubConfig.Exchange(context.Background(), code)
	utils.HandleError("Code-taken exchange failed", err)

	client := githubConfig.Client(context.Background(), token)
	resp, err := client.Get(GithubAuthURL)
	// fmt.Println(resp.Body)
	utils.HandleError("Failed to fetch user data from github:", err)

	defer resp.Body.Close()

	// parse response
	value := ParseOauthResponse(resp)
	fmt.Println(value)

	account := OauthAccount{
		Name:  utils.AssertString(value["name"]),
		Email: utils.AssertString(value["email"]),
	}

	LoginUserWithOauth(w, r, tmpl, account)
}
