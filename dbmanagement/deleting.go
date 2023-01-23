package dbmanagement

import (
	"database/sql"
	"fmt"
	"forum/utils"
	"log"
)

func DeleteFromTableWithUUID(table string, UUID string) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	log.Println("Deleting "+table+" record for uuid: ", UUID, "...")

	deleteRowStatement := "DELETE FROM " + table + " WHERE uuid = ?"
	statement, err := db.Prepare(deleteRowStatement)
	utils.HandleError("Delete Prepare failed: ", err)

	_, err = statement.Exec(UUID)
	utils.HandleError("Statement Exec failed: ", err)
}

func DeleteAllPostsWithTag(tagName string) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	log.Println("Deleting all post records for posts with tag: ", tagName, "...")

	listOfPostsToDelete := SelectAllPostsFromTag(tagName)

	fmt.Println("trying to delete the posts: ", listOfPostsToDelete)

	for _, post := range listOfPostsToDelete {
		deleteRowStatement := "DELETE FROM Posts WHERE uuid = ?"
		statement, err := db.Prepare(deleteRowStatement)
		utils.HandleError("Delete Prepare failed: ", err)

		_, err = statement.Exec(post.UUID)
		utils.HandleError("Statement Exec failed: ", err)
	}
}
