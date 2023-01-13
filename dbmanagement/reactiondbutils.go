package dbmanagement

import (
	"database/sql"
	"forum/utils"
	"log"
)

// add either a like or dislike or return to neutral with reaction taking values of 1,-1 or 0
/*
 */
func AddReactionToPost(userId string, postId string, reaction int) {
	if reaction > 1 || reaction < -1 || reaction == 0 {
		log.Println("Incorrect reaction integer")
		return
	}

	var reactionStatement string
	var postUpdateStatement string

	switch SelectReactionFromPost(userId, postId) {
	case -2:
		reactionStatement = `
		INSERT OR IGNORE INTO ReactedPosts(userId, postId, reaction) 
		VALUES (?, ?, ?)
	`
		if reaction == 1 {
			postUpdateStatement = `
	UPDATE Posts 
	SET likes = likes + 1
	WHERE uuid = ?
`
		} else if reaction == -1 {
			postUpdateStatement = `
	UPDATE Posts 
	SET dislikes = dislikes + 1 
	WHERE uuid = ?
`
		}
	case -1:
		if reaction == 1 {
			reactionStatement = `
		UPDATE ReactedPosts 
		SET reaction = ? 
		WHERE userId = ? and postId = ?
	`
			postUpdateStatement = `
		UPDATE Posts 
		SET likes = likes + 1, dislikes = dislikes - 1 
		WHERE uuid = ?
	`
		} else if reaction == -1 {
			reactionStatement = `
		UPDATE ReactedPosts 
		SET reaction = ? + 1
		WHERE userId = ? and postId = ?
	`
			postUpdateStatement = `
		UPDATE Posts 
		SET dislikes = dislikes - 1 
		WHERE uuid = ?
	`
		}
	case 0:
		reactionStatement = `
		UPDATE ReactedPosts 
		SET reaction = ? 
		WHERE userId = ? and postId = ?
	`
		if reaction == 1 {
			postUpdateStatement = `
		UPDATE Posts 
		SET likes = likes + 1
		WHERE uuid = ?
	`
		} else if reaction == -1 {
			postUpdateStatement = `
		UPDATE Posts 
		SET dislikes = dislikes + 1 
		WHERE uuid = ?
	`
		}
	case 1:
		if reaction == -1 {
			reactionStatement = `
		UPDATE ReactedPosts 
		SET reaction = ? 
		WHERE userId = ? and postId = ?
	`
			postUpdateStatement = `
		UPDATE Posts 
		SET likes = likes - 1, dislikes = dislikes + 1 
		WHERE uuid = ?
	`
		} else if reaction == 1 {
			reactionStatement = `
		UPDATE ReactedPosts 
		SET reaction = ? - 1 
		WHERE userId = ? and postId = ?
	`
			postUpdateStatement = `
		UPDATE Posts 
		SET likes = likes - 1 
		WHERE uuid = ?
	`
		}
	default:

	}

	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	statement, err := db.Prepare(reactionStatement)
	utils.HandleError("Like Prepare failed: ", err)

	ps, err := db.Prepare(postUpdateStatement)
	utils.HandleError("Post Update Prepare failed: ", err)

	if SelectReactionFromPost(userId, postId) == -2 {
		_, err = statement.Exec(userId, postId, reaction)
		utils.HandleError("statement Exec failed: ", err)
		_, err = ps.Exec(postId)
		utils.HandleError("post update Exec failed: ", err)
	} else {
		_, err = statement.Exec(reaction, userId, postId)
		utils.HandleError("statement Exec failed: ", err)
		_, err = ps.Exec(postId)
		utils.HandleError("post update Exec failed: ", err)
	}
}

func SelectReactionFromPost(postid, userid string) int {
	//-2 because you don't want to insert more entries into reaction table
	var reaction = -2
	var user string
	var post string
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	stm, err := db.Prepare("SELECT * FROM ReactedPosts WHERE userid = ? and postid = ?")
	utils.HandleError("Statement failed: ", err)

	err = stm.QueryRow(postid, userid).Scan(&user, &post, &reaction)
	utils.HandleError("Query Row failed: ", err)

	return reaction
}
