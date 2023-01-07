package rashdb

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
	content  string
	ownerId  string
	likes    int
	dislikes int
	time     time.Time
}

type Comment struct {
	UUID     string
	content  string
	postId   string
	ownerId  string
	likes    int
	dislikes int
	time     time.Time
}

var createUserTableStatement = `
	CREATE TABLE Users (
		uuid TEXT PRIMARY KEY,		
		name TEXT,
		email TEXT,
		password TEXT,
		permission TEXT
	)
;`

var createPostTableStatement = `
	CREATE TABLE Posts (
		uuid TEXT PRIMARY KEY,
		content TEXT,		
		ownerId TEXT,
		likes INTEGER,
		dislikes INTEGER,
		time DATETIME,
		FOREIGN KEY(ownerId) REFERENCES Users(uuid)
	)
;`

var createCommentTableStatement = `
	CREATE TABLE Comments (
		uuid TEXT PRIMARY KEY,
		content TEXT,
		postId TEXT,
		ownerId TEXT,
		likes INTEGER,
		dislikes INTEGER,
		time DATETIME,
		FOREIGN KEY(postId) REFERENCES Posts(uuid),
		FOREIGN KEY(ownerId) REFERENCES Users(uuid)
	)
;`

var createTagsTableStatement = `
	CREATE TABLE Tags (
		uuid TEXT PRIMARY KEY,
		tagname TEXT
	)
;`

var createTaggedPostsStatement = `
CREATE TABLE TaggedPosts (
		uuid TEXT PRIMARY KEY,
		tagId TEXT,
		postId TEXT,
		FOREIGN KEY(tagId) REFERENCES Tags(uuid),
		FOREIGN KEY(postId) REFERENCES Posts(uuid)
	)
;`

var createLikedPostsTableStatement = `
	CREATE TABLE LikedPosts (
		uuid TEXT PRIMARY KEY,
		userId TEXT,
		postId TEXT,
		liked INTEGER,
		disliked INTEGER,
		FOREIGN KEY(userId) REFERENCES Users(uuid),
		FOREIGN KEY(postId) REFERENCES Posts(uuid)
	)
;`

var createLikedCommentsTableStatement = `
	CREATE TABLE LikedComments (
		uuid TEXT PRIMARY KEY,
		userId TEXT,
		commentId TEXT,
		FOREIGN KEY(userId) REFERENCES Users(uuid),
		FOREIGN KEY(commentId) REFERENCES Comments(uuid)
	)
;`

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
	utils.HandleError("", err)
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
func InsertPost(content string, ownerId string, likes int, dislikes int, time time.Time) Post {
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
		log.Println("Post: ", UUID, " ", content, " ", ownerId, " ", likes, " ", dislikes, " ", time)
	}
}

func SelectUniqueUser(userName string) User {
	var user User
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	stm, err := db.Prepare("SELECT * FROM Users WHERE name = ?")
	utils.HandleError("Statement failed: ", err)

	err = stm.QueryRow(userName).Scan(&user.UUID, &user.Name, &user.Email, &user.Password, &user.Permission)
	utils.HandleError("Query Row failed: ", err)

	return user
}
