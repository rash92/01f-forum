package dbmanagement

import (
	"database/sql"
	"forum/utils"
	"log"
)

func AddNotification(receivingUserId, postId, commentId, sendingUserId string, reaction int) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	log.Println("Inserting notification record...")

	UUID := GenerateUUIDString()
	insertNotificationData := "INSERT INTO Notifications(UUID, receivingUserId, postId, commentId, sendingUserId, reaction) VALUES (?, ?, ?, ?, ?, ?)"
	statement, err := db.Prepare(insertNotificationData)
	utils.HandleError("User Prepare failed: ", err)

	_, err = statement.Exec(UUID, receivingUserId, postId, commentId, sendingUserId, reaction)
	utils.HandleError("Statement Exec failed: ", err)
}

func SelectAllNotificationsFromUser(receiver string) []Notification {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	row, err := db.Query("SELECT * FROM Notifications WHERE receivingUserId = ?", receiver)
	utils.HandleError("Post from User query failed: ", err)
	defer row.Close()

	var allNotifications []Notification

	for row.Next() {
		var currentNotification Notification
		row.Scan(&currentNotification.UUID, &currentNotification.Receiver, &currentNotification.PostId, &currentNotification.CommentId, &currentNotification.Sender, &currentNotification.Reaction)
		allNotifications = append(allNotifications, currentNotification)
	}
	return allNotifications
}
