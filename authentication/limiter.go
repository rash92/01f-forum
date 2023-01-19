package auth

import (
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

	if time.Now().After(endTime) {
		user.LimitTokens = 10
	}

	return user

}
