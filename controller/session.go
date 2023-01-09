package controller

import (
	"forum/utils"
	"net/http"
)

/*
Returns the cookie value of the current session that gives a sessions ID.  Used to determine which user is using the program.
*/
func Session(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie("_cookie")
	utils.HandleError("cookie err:", err)
	value := cookie.Value
	return value, err
}
