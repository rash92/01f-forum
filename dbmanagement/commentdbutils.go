package dbmanagement

import (
	"database/sql"
	"forum/utils"
	"log"
	"time"
)

/*
Inserts post into database with the relevant data, likes and dislikes should be set to 0 for most cases.  Each comment has it's own UUID.
*/
func InsertComment(content string, postId string, ownerId string, likes int, dislikes int, time time.Time) Comment {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	log.Println("Inserting comment record...")

	UUID := GenerateUUIDString()
	insertCommentData := "INSERT INTO Comments(UUID, content, postId, ownerId, likes, dislikes, time) VALUES (?, ?, ?, ?, ?, ?, ?)"
	statement, err := db.Prepare(insertCommentData)
	utils.HandleError("User Prepare failed: ", err)

	_, err = statement.Exec(UUID, content, postId, ownerId, likes, dislikes, time)
	utils.HandleError("Statement Exec failed: ", err)

	return Comment{UUID, content, postId, ownerId, likes, dislikes, time}
}

/*
Displays all comments from the database in the console.  Only for internal use.
*/
func DisplayAllComments() {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	row, err := db.Query("SELECT * FROM Comments ORDER BY time")
	utils.HandleError("User query failed: ", err)
	defer row.Close()

	for row.Next() {
		var UUID string
		var content string
		var postId string
		var ownerId string
		var likes int
		var dislikes int
		var time time.Time
		row.Scan(&UUID, &content, &postId, &ownerId, &likes, &dislikes, &time)
		owner := SelectUserFromUUID(ownerId)
		log.Println("Comment: ", UUID, " replying to: ", postId, " content: ", content, " owner: ", owner.Name, " likes ", likes, " dislikes ", dislikes, " time ", time)
	}
}

/*
Select and return Comment by it's UUID.
*/
func SelectCommentFromUUID(UUID string) Comment {
	var comment Comment
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	stm, err := db.Prepare("SELECT * FROM Comments WHERE uuid = ?")
	utils.HandleError("Statement failed: ", err)

	err = stm.QueryRow(UUID).Scan(&comment.UUID, &comment.Content, &comment.PostId, &comment.OwnerId, &comment.Likes, &comment.Dislikes, &comment.Time)
	utils.HandleError("Query Row failed: ", err)

	return comment
}

/*
Finds and returns all the comments made by a particular User using the ownerID (User's UUID).
*/
func SelectAllCommentsFromUser(ownerId string) []Comment {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	row, err := db.Query("SELECT * FROM Comments WHERE ownerId = ?", ownerId)
	utils.HandleError("User query failed: ", err)
	defer row.Close()

	var allComments []Comment

	for row.Next() {
		var currentComment Comment
		row.Scan(&currentComment.UUID, &currentComment.Content, &currentComment.PostId, &currentComment.OwnerId, &currentComment.Likes, &currentComment.Dislikes, &currentComment.Time)
		allComments = append(allComments, currentComment)
	}
	return allComments
}

/*
Finds and returns all comments on a particular post using a post's UUID.  Used for displaying all comments on the website, more specifically when you are querying a particular post.
*/
func SelectAllCommentsFromPost(postId string) []Comment {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	row, err := db.Query("SELECT * FROM Comments WHERE postId = ?", postId)
	utils.HandleError("User query failed: ", err)
	defer row.Close()

	var allComments []Comment

	for row.Next() {
		var currentComment Comment
		row.Scan(&currentComment.UUID, &currentComment.Content, &currentComment.PostId, &currentComment.OwnerId, &currentComment.Likes, &currentComment.Dislikes, &currentComment.Time)
		allComments = append(allComments, currentComment)
	}
	return allComments
}
