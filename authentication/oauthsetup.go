package auth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

func GoogleSetupConfig() *oauth2.Config {
	// Your credentials should be obtained from the Google
	// Developer Console (https://console.developers.google.com).
	conf := &oauth2.Config{
		ClientID:     "518245388319-qj8tmfb4ue1hfodjophsp9bprfe0om66.apps.googleusercontent.com",
		ClientSecret: "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",
		Endpoint:     google.Endpoint,
		RedirectURL:  "https://localhost:8080/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}
	return conf
}

func GithubSetupConfig() *oauth2.Config {
	// Oauth configuration for github
	conf := &oauth2.Config{
		ClientID:     "7a5de3b35a748e59ec9b",
		ClientSecret: "5256fcbe219c58b572029d7443dc3a996c7d378a",
		Endpoint:     github.Endpoint,
		RedirectURL:  "https://localhost:8080/github/callback",
		Scopes: []string{
			"user",
			"user:email",
		},
	}
	return conf
}

func FacebookSetupConfig() *oauth2.Config {
	// Oauth configuration for github
	conf := &oauth2.Config{
		ClientID:     "481977910781218",
		ClientSecret: "a23ef81bcd37b0c232c8eeb964d1402d",
		Endpoint:     facebook.Endpoint,
		RedirectURL:  "https://localhost:8080/facebook/callback",
		Scopes: []string{
			"user",
		},
	}
	return conf
}
