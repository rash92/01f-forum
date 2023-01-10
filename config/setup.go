package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func SetupConfig() *oauth2.Config {
	// Your credentials should be obtained from the Google
	// Developer Console (https://console.developers.google.com).
	conf := &oauth2.Config{
		ClientID:     "518245388319-qj8tmfb4ue1hfodjophsp9bprfe0om66.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "https://localhost:8080/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return conf
}
