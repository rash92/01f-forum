package dbmanagement

import (
	"database/sql"
	"fmt"
	"forum/utils"
	"log"
	"strings"
	"time"
)

/*
Inserts post into database with the relevant data, likes and dislikes should be set to 0 for most cases.  Each post has it's own UUID.
*/
func InsertPost(title string, content string, ownerId string, likes int, dislikes int, inputtime time.Time) Post {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	log.Println("Inserting post record...")

	UUID := GenerateUUIDString()
	insertPostData := "INSERT INTO Posts(UUID, title, content, ownerId, likes, dislikes, time) VALUES (?, ?, ?, ?, ?, ?, ?)"
	statement, err := db.Prepare(insertPostData)
	utils.HandleError("User Prepare failed: ", err)

	_, err = statement.Exec(UUID, title, content, ownerId, likes, dislikes, inputtime)
	utils.HandleError("Statement Exec failed: ", err)

	tags := SelectAllTagsFromPost(UUID)

	return Post{UUID, title, content, ownerId, likes, dislikes, tags, inputtime, strings.TrimSuffix(inputtime.Format(time.RFC822), "UTC"), 0}
}

/*
Displays all posts from the database in the console.  Only for internal use.
*/
func DisplayAllPosts() {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	row, err := db.Query("SELECT * FROM Posts ORDER BY time")
	utils.HandleError("Display posts query failed: ", err)
	defer row.Close()

	for row.Next() {
		var UUID string
		var title string
		var content string
		var ownerId string
		var likes int
		var dislikes int
		var time time.Time
		row.Scan(&UUID, &title, &content, &ownerId, &likes, &dislikes, &time)
		owner, err := SelectUserFromUUID(ownerId)
		utils.HandleError("unable to get user to display post", err)
		log.Println("Post: ", UUID, " content: ", content, " owner: ", owner.Name, " likes ", likes, " dislikes ", dislikes, " time ", time, "tags ", SelectAllTagsFromPost(UUID))
	}
}

/*
Finds a specific post based on the UUID (of the post).  Could be used for when bringing up a particular post onto one page.
*/
func SelectPostFromUUID(UUID string) Post {
	var post Post
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	stm, err := db.Prepare("SELECT * FROM Posts WHERE uuid = ?")
	utils.HandleError("Statement failed: ", err)

	err = stm.QueryRow(UUID).Scan(&post.UUID, &post.Title, &post.Content, &post.OwnerId, &post.Likes, &post.Dislikes, &post.Time)
	post.FormattedTime = strings.TrimSuffix(post.Time.Format(time.RFC822), "UTC")
	post.NumOfComments = len(SelectAllCommentsFromPost(post.UUID))
	post.Tags = SelectAllTagsFromPost(post.UUID)
	utils.HandleError("Query Row failed: ", err)

	return post
}

/*
Gathers all the posts from the database and returns them as an array of Post struct.  This function is used when displaying all the posts on the forum website.
*/
func SelectAllPosts() []Post {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	row, err := db.Query("SELECT * FROM Posts")
	utils.HandleError("All posts query failed: ", err)
	defer row.Close()

	var allPosts []Post
	for row.Next() {
		var currentPost Post
		row.Scan(&currentPost.UUID, &currentPost.Title, &currentPost.Content, &currentPost.OwnerId, &currentPost.Likes, &currentPost.Dislikes, &currentPost.Time)
		currentPost.FormattedTime = strings.TrimSuffix(currentPost.Time.Format(time.RFC822), "UTC")
		currentPost.NumOfComments = len(SelectAllCommentsFromPost(currentPost.UUID))
		currentPost.Tags = SelectAllTagsFromPost(currentPost.UUID)
		allPosts = append(allPosts, currentPost)
	}
	return allPosts
}

/*
Similar to SelectAllPosts() but for a specific user.  Uses the ownerID (users UUID) to specify which user and returns all the posts created by that user.
*/
func SelectAllPostsFromUser(ownerId string) []Post {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	row, err := db.Query("SELECT * FROM Posts WHERE ownerId = ?", ownerId)
	utils.HandleError("Post from User query failed: ", err)
	defer row.Close()

	var allPosts []Post

	for row.Next() {
		var currentPost Post
		row.Scan(&currentPost.UUID, &currentPost.Title, &currentPost.Content, &currentPost.OwnerId, &currentPost.Likes, &currentPost.Dislikes, &currentPost.Time)
		currentPost.FormattedTime = strings.TrimSuffix(currentPost.Time.Format(time.RFC822), "UTC")
		currentPost.NumOfComments = len(SelectAllCommentsFromPost(currentPost.UUID))
		currentPost.Tags = SelectAllTagsFromPost(currentPost.UUID)
		allPosts = append(allPosts, currentPost)
	}
	return allPosts
}

func SelectAllLikedPostsFromUser(user User) []Post {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	row, err := db.Query("SELECT * FROM Posts WHERE ownerId = ?", user.Name)
	utils.HandleError("Post from User query failed: ", err)
	defer row.Close()

	var allPosts []Post

	for row.Next() {
		var currentPost Post
		row.Scan(&currentPost.UUID, &currentPost.Title, &currentPost.Content, &currentPost.OwnerId, &currentPost.Likes, &currentPost.Dislikes, &currentPost.Time)
		currentPost.FormattedTime = strings.TrimSuffix(currentPost.Time.Format(time.RFC822), "UTC")
		currentPost.NumOfComments = len(SelectAllCommentsFromPost(currentPost.UUID))
		currentPost.Tags = SelectAllTagsFromPost(currentPost.UUID)
		if SelectReactionFromPost(currentPost.UUID, user.UUID) == 1 {
			allPosts = append(allPosts, currentPost)
		}
	}
	return allPosts
}

/*
Similar to SelectAllPosts() but for a specific user.  Uses the ownerID (users UUID) to specify which user and returns all the posts created by that user.
*/
func SelectAllPostsFromTag(tagName string) []Post {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	tag, err := SelectTagFromName(tagName)
	utils.HandleError("couldn't find tag", err)

	row, err := db.Query("SELECT postId FROM TaggedPosts WHERE tagId = ?", tag.UUID)
	utils.HandleError("Tag query failed: ", err)
	defer row.Close()

	var allPosts []Post

	for row.Next() {
		var currentPostId string
		var currentPost Post
		row.Scan(&currentPostId)
		currentPost = SelectPostFromUUID(currentPostId)
		currentPost.FormattedTime = strings.TrimSuffix(currentPost.Time.Format(time.RFC822), "UTC")
		currentPost.Tags = SelectAllTagsFromPost(currentPost.UUID)
		currentPost.NumOfComments = len(SelectAllCommentsFromPost(currentPost.UUID))
		fmt.Println("found post from tag: ", tagName, "the post: ", currentPost)
		allPosts = append(allPosts, currentPost)
	}

	return allPosts
}
