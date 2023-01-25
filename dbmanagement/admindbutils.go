package dbmanagement

import (
	"database/sql"
	"forum/utils"
)

func CreateAdminRequest(RequestFromId string, RequestFromName string, content string) AdminRequest {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	utils.WriteMessageToLogFile("Inserting admin request record...")

	UUID := GenerateUUIDString()
	insertAdminRequestData := "INSERT INTO AdminRequests(UUID, requestfromid, requestfromname, content) VALUES (?, ?, ?, ?)"
	statement, err := db.Prepare(insertAdminRequestData)
	utils.HandleError("User Prepare failed: ", err)

	utils.WriteMessageToLogFile("admint request content is: " + content)

	_, err = statement.Exec(UUID, RequestFromId, RequestFromName, content)
	utils.HandleError("Statement Exec failed: ", err)

	return AdminRequest{UUID, RequestFromId, RequestFromName, content}
}

func SelectAllAdminRequests() []AdminRequest {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	row, err := db.Query("SELECT * FROM AdminRequests")
	utils.HandleError("Admin Request query failed: ", err)
	defer row.Close()

	var allAdminRequests []AdminRequest
	for row.Next() {
		var currentAdminRequest AdminRequest
		row.Scan(&currentAdminRequest.UUID, &currentAdminRequest.RequestFromId, &currentAdminRequest.RequestFromName, &currentAdminRequest.Content)
		allAdminRequests = append(allAdminRequests, currentAdminRequest)
	}
	return allAdminRequests
}
