package dbmanagement

import (
	"database/sql"
	"fmt"
	"forum/utils"
	"os"
)

func DeleteFromTableWithUUID(table string, UUID string) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	utils.WriteMessageToLogFile("Deleting " + table + " record for uuid: " + UUID + "...")

	deleteRowStatement := "DELETE FROM " + table + " WHERE uuid = ?"
	statement, err := db.Prepare(deleteRowStatement)
	utils.HandleError("Delete Prepare failed: ", err)

	_, err = statement.Exec(UUID)
	utils.HandleError("Statement Exec failed: ", err)
}

func DeleteFromTableWithPostId(table string, postId string) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	utils.WriteMessageToLogFile("Deleting " + table + " record for postId: " + postId + "...")

	deleteRowStatement := "DELETE FROM " + table + " WHERE postId = ?"
	statement, err := db.Prepare(deleteRowStatement)
	utils.HandleError("Delete Prepare failed: ", err)

	_, err = statement.Exec(postId)
	utils.HandleError("Statement Exec failed: ", err)
}

func DeletePostWithUUID(UUID string) {
	post := SelectPostFromUUID(UUID)
	if post.ImageName != "" {
		os.Remove("." + post.ImageName)
	}
	comments := SelectAllCommentsFromPost(UUID)
	for _, comment := range comments {
		DeleteFromTableWithUUID("comments", comment.UUID)
	}
	notifications := SelectAllNotificationsFromUUID(UUID)
	for _, notification := range notifications {
		DeleteFromTableWithUUID("Notifications", notification.PostId)
	}
	DeleteFromTableWithUUID("posts", UUID)
}

func DeleteAllPostsWithTag(tagName string) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	utils.WriteMessageToLogFile("Deleting all post records for posts with tag: " + tagName + "...")

	listOfPostsToDelete := SelectAllPostsFromTag(tagName)

	message := fmt.Sprint("trying to delete the posts: ", listOfPostsToDelete)
	utils.WriteMessageToLogFile(message)

	for _, post := range listOfPostsToDelete {
		DeletePostWithUUID(post.UUID)
	}
}

func DeleteUser(name string) error {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	statement := "DELETE FROM Users WHERE name = ?"
	stm, err := db.Prepare(statement)
	utils.HandleError("Failed to delete user statement in", err)

	defer stm.Close()

	res, err := stm.Exec(name)
	utils.HandleError("Failed to delete user in", err)

	n, err := res.RowsAffected()
	utils.HandleError("Rows affected error in", err)

	message := fmt.Sprintf("The statement has affected %d rows\n", n)
	utils.WriteMessageToLogFile(message)
	return err
}

// Delete session from database
func DeleteSessionByUUID(UUID string) (err error) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	statement := "DELETE FROM Sessions WHERE uuid = ?"
	stm, err := db.Prepare(statement)
	utils.HandleError("Failed to delete session by uuid:", err)

	defer stm.Close()

	res, err := stm.Exec(UUID)

	n, err := res.RowsAffected()
	utils.HandleError("Rows affected error:", err)

	message := fmt.Sprintf("The statement has affected %d rows", n)
	utils.WriteMessageToLogFile(message)
	return
}

// Delete all session from database
func DeleteAllSessions() (err error) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	res, err := db.Exec("DELETE FROM Sessions")
	utils.HandleError("Failed to delete all sessions:", err)

	n, err := res.RowsAffected()
	utils.HandleError("Rows affected error:", err)

	message := fmt.Sprintf("The statement has affected %d rows", n)
	utils.WriteMessageToLogFile(message)
	return err
}
