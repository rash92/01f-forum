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

	//if the user has not logged in, they are still able to view the forum page.
	//There will be no cookie value if they have not logged in
	var value string
	if err == nil {
		value = cookie.Value
	}
	return value, err
}
