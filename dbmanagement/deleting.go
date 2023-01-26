package dbmanagement

import (
	"database/sql"
	"fmt"
	"forum/utils"
	"log"
	"os"
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

func DeleteFromTableWithPostId(table string, postId string) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	log.Println("Deleting "+table+" record for postId: ", postId, "...")

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
	DeleteFromTableWithUUID("posts", UUID)
}

func DeleteAllPostsWithTag(tagName string) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	log.Println("Deleting all post records for posts with tag: ", tagName, "...")

	listOfPostsToDelete := SelectAllPostsFromTag(tagName)

	fmt.Println("trying to delete the posts: ", listOfPostsToDelete)

	for _, post := range listOfPostsToDelete {
		DeletePostWithUUID(post.UUID)
	}
}
