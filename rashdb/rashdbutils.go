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
	UUID       uuid.UUID
	Name       string
	Email      string
	Password   string
	Permission string
}

type Post struct {
	UUID     uuid.UUID
	content  string
	owner    *User
	likes    int
	dislikes int
	time     time.Time
}

type Comment struct {
	UUID       uuid.UUID
	content    string
	replyingto *Post
	owner      *User
	likes      int
	dislikes   int
	time       time.Time
}

var createUserTableDB = `
	CREATE TABLE Users (
		uuid TEXT PRIMARY KEY,		
		name TEXT,
		email TEXT,
		password TEXT,
		permission TEXT
	)
;`

var createPostTableDB = `
	CREATE TABLE Posts (
		uuid TEXT PRIMARY KEY,
		content TEXT,		
		owner TEXT,
		likes INTEGER,
		dislikes INTEGER,
		time DATETIME,
		FOREIGN KEY(owner) REFERENCES Users(uuid)
	)
;`

var createCommentTableDB = `
	CREATE TABLE Comments (
		uuid TEXT PRIMARY KEY,
		content TEXT,
		replyingto TEXT,
		owner TEXT,
		likes INTEGER,
		dislikes INTEGER,
		time DATETIME,
		FOREIGN KEY(replyingto) REFERENCES Posts(uuid),
		FOREIGN KEY(owner) REFERENCES Users(uuid)
	)
;`

var createTagsTableDB = `
	CREATE TABLE Tags (
		uuid TEXT PRIMARY KEY,
		tagname TEXT
	)
;`

var createTaggedPostsDB = `
CREATE TABLE TaggedPosts (
		uuid TEXT PRIMARY KEY,
		tag TEXT,
		post TEXT,
		FOREIGN KEY(tag) REFERENCES Tags(uuid),
		FOREIGN KEY(post) REFERENCES Posts(uuid)
	)
;`

var createLikedPostsTableDB = `
	CREATE TABLE LikedPosts (
		uuid TEXT PRIMARY KEY,
		user TEXT,
		post TEXT,
		FOREIGN KEY(user) REFERENCES Users(uuid),
		FOREIGN KEY(post) REFERENCES Posts(uuid)
	)
;`

var createLikedCommentsTableDB = `
	CREATE TABLE LikedComments (
		uuid TEXT PRIMARY KEY,
		user TEXT,
		comment TEXT,
		FOREIGN KEY(user) REFERENCES Users(uuid),
		FOREIGN KEY(comment) REFERENCES Comments(uuid)
	)
;`

func CreateDatabaseWithTables() {
	forumDB := CreateDatabase()
	defer forumDB.Close()

	CreateTable(forumDB, createUserTableDB)
	CreateTable(forumDB, createPostTableDB)

	log.Println("forum.db created successfully!")
}

func CreateDatabase() *sql.DB {
	// os.Remove("forum.db")
	log.Println("Creating forum.db...")
	file, err := os.Create("forum.db")
	utils.HandleError("", err)

	file.Close()

	forumDB, err := sql.Open("sqlite3", "./forum.db?_foreign_keys=on")
	utils.HandleError("", err)

	return forumDB
}

func CreateTable(db *sql.DB, table string) {
	statement, err := db.Prepare(table)
	utils.HandleError("", err)
	statement.Exec()
}

func InsertUser(UUID string, name string, email string, password string, permission string) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	log.Println("Inserting user record...")

	insertUserData := "INSERT INTO Users(UUID, name, email, password, permission) VALUES (?, ?, ?, ?, ?)"
	statement, err := db.Prepare(insertUserData)
	utils.HandleError("User Prepare failed: ", err)

	_, err = statement.Exec(UUID, name, email, password, permission)
	utils.HandleError("Statement Exec failed: ", err)
}

func DisplayAllUsers() {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	row, err := db.Query("SELECT * FROM Users ORDER BY name")
	utils.HandleError("User query failed: ", err)
	defer row.Close()

	for row.Next() {
		var user_id int
		var UUID string
		var name string
		var email string
		var password string
		var permission string
		row.Scan(&user_id, &UUID, &name, &email, &password, &permission)
		log.Println("User: ", user_id, " ", UUID, " ", name, " ", email, " ", password, " ", permission)
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
