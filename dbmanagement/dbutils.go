package dbmanagement

import (
	"database/sql"
	"forum/utils"
	"log"
	"os"
	"time"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
)

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
		tag TEXT,
		time DATETIME,
		FOREIGN KEY (ownerId) REFERENCES Users(uuid)
		FOREIGN KEY (tag) REFERENCES Tags(tagname)
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
var createReactionPostsTableStatement = `
	CREATE TABLE ReactedPosts (
		userId TEXT NOT NULL,
		postId TEXT NOT NULL,
		reaction INTEGER,
		FOREIGN KEY (userId) REFERENCES Users(uuid),
		FOREIGN KEY (postId) REFERENCES Posts(uuid),
		PRIMARY KEY (userId, postId)
	);`

var createLikedCommentsTableStatement = `
	CREATE TABLE LikedComments  (
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

var createAdminRequestTableStatement = `
	CREATE TABLE AdminRequests (
		uuid TEXT NOT NULL PRIMARY KEY,
		requestfromid TEXT,
		requestfromname TEXT,
		content TEXT,
		FOREIGN KEY (requestfromid) REFERENCES Users(uuid),
		FOREIGN KEY (requestfromid) REFERENCES Users(name)	
	);`

/*
Only used to create brand new databases, wiping all previous data in the process.
To be used when initially implementing database or clearing data after testing.
*/
func CreateDatabaseWithTables() {
	forumDB := CreateDatabase("forum")
	defer forumDB.Close()

	CreateTable(forumDB, createUserTableStatement)
	CreateTable(forumDB, createPostTableStatement)
	CreateTable(forumDB, createCommentTableStatement)
	CreateTable(forumDB, createTagsTableStatement)
	CreateTable(forumDB, createTaggedPostsStatement)
	CreateTable(forumDB, createReactionPostsTableStatement)
	CreateTable(forumDB, createLikedCommentsTableStatement)
	CreateTable(forumDB, createSessionTableStatement)
	CreateTable(forumDB, createAdminRequestTableStatement)

	log.Println("forum.db created successfully!")
}

/*
Creates a new database file to store tables.  If database already exists, it is REMOVED.  Beware of losing data.
*/
func CreateDatabase(name string) *sql.DB {
	os.Remove(name + ".db")
	log.Println("Creating " + name + ".db...")
	file, err := os.Create(name + ".db")
	utils.HandleError("", err)

	file.Close()

	forumDB, err := sql.Open("sqlite3", "./"+name+".db?_foreign_keys=on")
	utils.HandleError("", err)

	return forumDB
}

/*
Creates a table within a specified database
*/
func CreateTable(db *sql.DB, table string) {
	statement, err := db.Prepare(table)
	utils.HandleError(table, err)
	statement.Exec()
}

/*
Generates a new UUID and returns a string of that new number.
*/
func GenerateUUIDString() string {
	UUID, err := uuid.NewV4()
	utils.HandleError("problem generating uuid", err)
	return UUID.String()
}

/*
Used to provide specific information for when a user logs in by cross referencing their userID.
Creates and returns a new session when the user successfully logs in to their account.
The sessions has its own UUID, contains the usersID (user's UUID), and the time that it was created.
*/
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
