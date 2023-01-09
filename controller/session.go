package controller

import (
	utils "forum/helpers"
	"net/http"
)

func Session(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie("_cookie")
	utils.HandleError("cookie err:", err)

	// log.Println("cookie:", cookie)

	// log.Println("cookie value:", cookie.Value)

	// cookieValue := strings.Split(cookie.Value, "=")

	// log.Println("cookie value:", cookieValue)

	value := cookie.Value

	// if err == nil {
	// 	sess = user.Session{Uuid: cookie.Value}
	// 	if ok, _ := sess.Check(); !ok {
	// 		err = errors.New("Invalid session")
	// 	}
	// }

	return string(value), err
}
