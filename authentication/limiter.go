package auth

import (
	"fmt"
	"forum/dbmanagement"
	"forum/utils"
	"net/http"
	"time"
)

const Limit = 3

func LimitRequests(w http.ResponseWriter, r *http.Request, user dbmanagement.User) dbmanagement.User {
	limitTime := time.Minute * 2
	userSession, err := user.ReturnSession(user.UUID)
	utils.HandleError("unable to get session for :", err)
	startTime := userSession.CreatedAt
	endTime := startTime.Add(limitTime)

	go func() {
		for {
			time.Sleep(15 * time.Second)
			user, err := dbmanagement.SelectUserFromSession(userSession.UUID)
			utils.HandleError("go routine problem :", err)
			if CheckTime(endTime) {
				dbmanagement.UpdateUserToken(userSession.UserId, Limit)
				startTime = endTime
				endTime = startTime.Add(limitTime)
			}
			fmt.Println(user.LimitTokens)
		}
	}()

	return user
}

func CheckTime(endTime time.Time) bool {
	return time.Now().After(endTime)
}
