package controller

import (
	"time"
)

type User struct {
	ID         int
	UUID       string
	Name       string
	Email      string
	Password   string
	Permission string
}

type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	CreatedAt time.Time
}

// func session(writer http.ResponseWriter, request *http.Request) (sess data.Session, err error) {
// 	cookie, err := request.Cookie("_cookie")

// 	if err == nil {
// 		sess = data.Session{Uuid: cookie.Value}
// 		if ok, _ := sess.Check(); !ok {
// 			err = errors.New("Invalid session")
// 		}
// 	}
// 	return
// }
