package dbmanagement

import (
	"database/sql"
	utils "forum/helpers"
	"log"
	"os"
	"time"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	UUID       string
	Name       string
	Email      string
	Password   string
	Permission string
}

type Post struct {
	UUID     string
	Content  string
	OwnerId  string
	Likes    int
	Dislikes int
	Time     time.Time
}

type Comment struct {
	UUID     string
	Content  string
	PostId   string
	OwnerId  string
	Likes    int
	Dislikes int
	Time     time.Time
}

type Session struct {
	UUID      string
	UserId    string
	CreatedAt time.Time
}

var createUserTableStatement = `
	CREATE TABLE Users (
		uuid TEXT NOT NULL PRIMARY KEY,		
		name TEXT UNIQUE,
		email TEXT,
		password TEXT,
		permission TEXT
	);`

var createPostTableStatement = `
	CREATE TABLE Posts (
		uuid TEXT NOT NULL PRIMARY KEY,
		content TEXT,		
		ownerId TEXT,
		likes INTEGER,
		dislikes INTEGER,
		time DATETIME,
		FOREIGN KEY (ownerId) REFERENCES Users(uuid)
	);`

var createCommentTableStatement = `
	CREATE TABLE Comments (
		uuid TEXT NOT NULL PRIMARY KEY,
		content TEXT,
		postId TEXT,
		ownerId TEXT,
		likes INTEGER,
		dislikes INTEGER,
		time DATETIME,
		FOREIGN KEY (postId) REFERENCES Posts(uuid),
		FOREIGN KEY (ownerId) REFERENCES Users(uuid)
	);`

var createTagsTableStatement = `
	CREATE TABLE Tags (
		uuid TEXT NOT NULL PRIMARY KEY,
		tagname TEXT
	);`

var createTaggedPostsStatement = `
CREATE TABLE TaggedPosts (
		tagId TEXT NOT NULL,
		postId TEXT NOT NULL,
		FOREIGN KEY (tagId) REFERENCES Tags(uuid),
		FOREIGN KEY (postId) REFERENCES Posts(uuid),
		PRIMARY KEY (tagId, postId)
	);`

// can represent like as 1, dislike as -1 and neither as 0 as a single value in the reaction field
var createLikedPostsTableStatement = `
	CREATE TABLE LikedPosts (
		userId TEXT NOT NULL,
		postId TEXT NOT NULL,
		reaction INTEGER,
		FOREIGN KEY (userId) REFERENCES Users(uuid),
		FOREIGN KEY (postId) REFERENCES Posts(uuid),
		PRIMARY KEY (userId, postId)
	);`

var createLikedCommentsTableStatement = `
	CREATE TABLE LikedComments (
		userId TEXT NOT NULL,
		commentId TEXT NOT NULL,
		reaction INTEGER,
		FOREIGN KEY (userId) REFERENCES Users(uuid),
		FOREIGN KEY (commentId) REFERENCES Comments(uuid),
		PRIMARY KEY (userId, commentId)
	);`

var createSessionTableStatement = `
	CREATE TABLE Sessions (
  uuid      TEXT NOT NULL PRIMARY KEY,
  userId    INTEGER REFERENCES Users(uuid),
  createdAt TIMESTAMP NOT NULL   
);`

func CreateDatabaseWithTables() {
	forumDB := CreateDatabase("forum")
	defer forumDB.Close()

	CreateTable(forumDB, createUserTableStatement)
	CreateTable(forumDB, createPostTableStatement)
	CreateTable(forumDB, createCommentTableStatement)
	CreateTable(forumDB, createTagsTableStatement)
	CreateTable(forumDB, createTaggedPostsStatement)
	CreateTable(forumDB, createLikedPostsTableStatement)
	CreateTable(forumDB, createLikedCommentsTableStatement)
	CreateTable(forumDB, createSessionTableStatement)

	log.Println("forum.db created successfully!")
}

func CreateDatabase(name string) *sql.DB {
	// os.Remove("forum.db")
	log.Println("Creating " + name + ".db...")
	file, err := os.Create(name + ".db")
	utils.HandleError("", err)

	file.Close()

	forumDB, err := sql.Open("sqlite3", "./"+name+".db?_foreign_keys=on")
	utils.HandleError("", err)

	return forumDB
}

func CreateTable(db *sql.DB, table string) {
	statement, err := db.Prepare(table)
	utils.HandleError(table, err)
	statement.Exec()
}

// uuid version that's allowed from packages outputs something of type UUID
// which needs to be converted to string to be used as we have been using the UUIDs from the other package
func GenerateUUIDString() string {
	UUID, err := uuid.NewV4()
	if err != nil {
		utils.HandleError("problem generating uuid", err)
	}
	return UUID.String()
}

// returning User in case we need to use the one we just created, if return value not assigned to anything then will still update database.
// This generates the UUID internally. Another option is to take a User as input instead of individual fields after generating uuid externally.
// similar situation for posts and comments below
func InsertUser(name string, email string, password string, permission string) User {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	log.Println("Inserting user record...")

	UUID := GenerateUUIDString()
	insertUserData := "INSERT INTO Users(UUID, name, email, password, permission) VALUES (?, ?, ?, ?, ?)"
	statement, err := db.Prepare(insertUserData)
	utils.HandleError("User Prepare failed: ", err)

	_, err = statement.Exec(UUID, name, email, password, permission)
	utils.HandleError("Statement Exec failed: ", err)

	return User{UUID, name, email, password, permission}
}

// there is the option to generate time internally rather than needing to pass it through (similarly for comments below) using time.Now
func InsertPost(content string, ownerId string, likes int, dislikes int, tag string, time time.Time) Post {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	log.Println("Inserting post record...")

	UUID := GenerateUUIDString()
	insertPostData := "INSERT INTO Posts(UUID, content, ownerId, likes, dislikes, time) VALUES (?, ?, ?, ?, ?, ?)"
	statement, err := db.Prepare(insertPostData)
	utils.HandleError("User Prepare failed: ", err)

	_, err = statement.Exec(UUID, content, ownerId, likes, dislikes, time)
	utils.HandleError("Statement Exec failed: ", err)

	return Post{UUID, content, ownerId, likes, dislikes, time}
}

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

// add either a like or dislike or return to neutral with reaction taking values of 1,-1 or 0
func AddReactionToPost(userId string, postId string, reaction int) {
	if reaction > 1 || reaction < -1 {
		return
	}
	addLikeStatement := `
		INSERT OR IGNORE INTO LikedPosts(userId, postId, reaction) 
		VALUES (?, ?, ?)
		UPDATE LikedPosts 
		SET reaction = ? 
		WHERE userId = ? and postId = ?
	`
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	statement, err := db.Prepare(addLikeStatement)
	utils.HandleError("Like Prepare failed: ", err)

	_, err = statement.Exec(userId, postId, reaction, reaction, userId, postId)
	utils.HandleError("statement Exec failed: ", err)

	// need to insert stuff to update the overall like count on the post itself
	// need to check if the user had already liked it or disliked it to tell how much to change it by
	// could also do a query on the entire likedposts and count up the likes and dislikes to update it
}

func DisplayAllUsers() {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	row, err := db.Query("SELECT * FROM Users ORDER BY name")
	utils.HandleError("User query failed: ", err)
	defer row.Close()

	for row.Next() {

		var UUID string
		var name string
		var email string
		var password string
		var permission string
		row.Scan(&UUID, &name, &email, &password, &permission)
		log.Println("User: ", UUID, " ", name, " ", email, " ", password, " ", permission)
	}
}

func DisplayAllPosts() {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	row, err := db.Query("SELECT * FROM Posts ORDER BY time")
	utils.HandleError("User query failed: ", err)
	defer row.Close()

	for row.Next() {
		var UUID string
		var content string
		var ownerId string
		var likes int
		var dislikes int
		var time time.Time
		row.Scan(&UUID, &content, &ownerId, &likes, &dislikes, &time)
		owner := SelectUserFromUUID(ownerId)
		log.Println("Post: ", UUID, " content: ", content, " owner: ", owner.Name, " likes ", likes, " dislikes ", dislikes, " time ", time)
	}
}

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

func SelectUserFromName(Name string) User {
	var user User
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	stm, err := db.Prepare("SELECT * FROM Users WHERE name = ?")
	utils.HandleError("Statement failed: ", err)

	err = stm.QueryRow(Name).Scan(&user.UUID, &user.Name, &user.Email, &user.Password, &user.Permission)
	utils.HandleError("Query Row failed: ", err)

	return user
}

func SelectUserFromEmail(Email string) User {
	var user User
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	stm, err := db.Prepare("SELECT * FROM Users WHERE email = ?")
	utils.HandleError("Statement failed: ", err)

	err = stm.QueryRow(Email).Scan(&user.UUID, &user.Name, &user.Email, &user.Password, &user.Permission)
	utils.HandleError("Query Row failed: ", err)

	return user
}

func SelectUserFromUUID(UUID string) User {
	var user User
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	stm, err := db.Prepare("SELECT * FROM Users WHERE uuid = ?")
	utils.HandleError("Statement failed: ", err)

	err = stm.QueryRow(UUID).Scan(&user.UUID, &user.Name, &user.Email, &user.Password, &user.Permission)
	utils.HandleError("Query Row failed: ", err)

	return user
}

func SelectPostFromUUID(UUID string) Post {
	var post Post
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	stm, err := db.Prepare("SELECT * FROM Posts WHERE uuid = ?")
	utils.HandleError("Statement failed: ", err)

	err = stm.QueryRow(UUID).Scan(&post.UUID, &post.Content, &post.OwnerId, &post.Likes, &post.Dislikes, &post.Time)
	utils.HandleError("Query Row failed: ", err)

	return post
}

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

func SelectAllPosts() []Post {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	row, err := db.Query("SELECT * FROM Posts")
	utils.HandleError("User query failed: ", err)
	defer row.Close()

	var allPosts []Post

	for row.Next() {
		var currentPost Post
		row.Scan(&currentPost.UUID, &currentPost.Content, &currentPost.OwnerId, &currentPost.Likes, &currentPost.Dislikes, &currentPost.Time)
		allPosts = append(allPosts, currentPost)
	}
	return allPosts
}

func SelectAllPostsFromUser(ownerId string) []Post {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	row, err := db.Query("SELECT * FROM Posts WHERE ownerId = ?", ownerId)
	utils.HandleError("User query failed: ", err)
	defer row.Close()

	var allPosts []Post

	for row.Next() {
		var currentPost Post
		row.Scan(&currentPost.UUID, &currentPost.Content, &currentPost.OwnerId, &currentPost.Likes, &currentPost.Dislikes, &currentPost.Time)
		allPosts = append(allPosts, currentPost)
	}
	return allPosts
}

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

func (user *User) CreateSession() (session Session, err error) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	statement := `INSERT INTO Sessions (uuid, userID, createdAt) values (?, ?, ?) returning uuid, userID, createdAt`

	stmt, err := db.Prepare(statement)
	utils.HandleError("session error:", err)

	defer stmt.Close()

	UUID := GenerateUUIDString()
	timeNow := time.Now()

	err = stmt.QueryRow(UUID, user.UUID, timeNow).Scan(&session.UUID, &session.UserId, &session.CreatedAt)
	return
}

// func (user *User) Session() (session Session, err error) {
// 	db, _ := sql.Open("sqlite3", "./forum.db")
// 	session = Session{}
// 	err = db.QueryRow("SELECT uuid, userID, createdAt FROM sessions WHERE userID = ?", user.UUID).
// 		Scan(&session.Uuid, &session.UserId, &session.CreatedAt)
// 	return
// }
