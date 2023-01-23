package auth

import (
	"fmt"
	"forum/dbmanagement"
	"forum/utils"
	"net/http"
	"time"
)

const Limit = 10

func LimitRequests(w http.ResponseWriter, r *http.Request, user dbmanagement.User) dbmanagement.User {
	limitTime := time.Minute * 5
	userSession, err := user.ReturnSession(user.UUID)
	utils.HandleError("unable to get session for :", err)
	startTime := userSession.CreatedAt
	endTime := startTime.Add(limitTime)

	go func(limitTime time.Duration, endTime time.Time) {
		for {
			time.Sleep(15 * time.Second)
			user, err := dbmanagement.SelectUserFromSession(userSession.UUID)
			utils.HandleError("go routine problem :", err)
			if CheckTime(limitTime, endTime) {
				dbmanagement.UpdateUserToken(userSession.UserId, Limit)
			}
			fmt.Println(user.LimitTokens)
		}
	}(limitTime, endTime)

	return user
}

func CheckTime(limitTime time.Duration, endTime time.Time) bool {
	return time.Now().After(endTime)
}
