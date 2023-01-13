package controller

import (
	"fmt"
	"forum/dbmanagement"
	"forum/utils"
	"net/http"
)

/*
Returns the cookie value of the current session that gives a sessions ID.  Used to determine which user is using the program.
*/
func GetSessionIDFromBrowser(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie("_cookie")
	utils.HandleError("Cannot get Cookie Err:", err)
	if err != nil {
		return "", err
	}
	value := cookie.Value
	return value, err
}

/*
Creates session that gives a sessions ID, used to determine which user is using the program.
*/
func CreateUserSessionCookie(w http.ResponseWriter, r *http.Request, user dbmanagement.User) error {
	session, err := user.CreateSession()
	utils.HandleError("Cannot create user session err:", err)
	cookie := http.Cookie{
		Name:     "_cookie",
		Value:    session.UUID,
		HttpOnly: true,
	}
	fmt.Println("google user cookie created here")
	http.SetCookie(w, &cookie)
	fmt.Println("google user cookie created here in the browser")
	return err
}
