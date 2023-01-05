package dbmanagement

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID         int
	UUID       string
	Name       string
	Email      string
	Password   string
	Permission string
}

type Post struct {
	Comment  string
	Userid   int
	Likes    int
	Dislikes int
	Time     interface{}
}

func CreateDatabase() {
	// os.Remove("forum.db")
	log.Println("Creating forum.db...")
	file, err := os.Create("forum.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	createUserTableDB := `
	CREATE TABLE Users (
		user_id INTEGER PRIMARY KEY AUTOINCREMENT,
		UUID TEXT NOT NULL,		
		name TEXT,
		email TEXT,
		password TEXT,
		permission TEXT
	  );`
	createPostTableDB := `
	CREATE TABLE Posts (
		post_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		comment TEXT,		
		user INTEGER,
		likes INTEGER,
		dislikes INTEGER,
		time DATETIME,
		tags TEXT,
		FOREIGN KEY(user) REFERENCES Users(user_id)
	  );`
	forumDB, _ := sql.Open("sqlite3", "./forum.db?_foreign_keys=on")
	defer forumDB.Close()
	CreateTable(forumDB, createUserTableDB)
	CreateTable(forumDB, createPostTableDB)
	log.Println("forum.db create successfully!")
}

func CreateTable(db *sql.DB, table string) {
	statement, err := db.Prepare(table)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
}

func InsertUser(UUID string, name string, email string, password string, permission string) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	log.Println("Inserting user record...")
	insertUserData := "INSERT INTO Users(UUID, name, email, password, permission) VALUES (?, ?, ?, ?, ?)"
	statement, err := db.Prepare(insertUserData)
	if err != nil {
		log.Fatalln("User Prepare failed: ", err.Error())
	}
	_, err = statement.Exec(UUID, name, email, password, permission)
	if err != nil {
		log.Fatalln("Statement Exec failed: ", err.Error())
	}
}

func InsertPost(text string, user int, time string, tags string) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	log.Println("Inserting post...")
	insertPostData := "INSERT INTO Posts(comment, user, time, tags) VALUES (?, ?, ?, ?)"
	statement, err := db.Prepare(insertPostData)
	if err != nil {
		log.Fatalln("Post Prepare failed: ", err.Error())
	}
	_, err = statement.Exec(text, user, time, tags)
	if err != nil {
		log.Fatalln("Statement Exec failed: ", err.Error())
	}
}

func DisplayAllUsers() {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	row, err := db.Query("SELECT * FROM Users ORDER BY name")
	if err != nil {
		log.Fatalln("User query failed: ", err.Error())
	}
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

func DisplayAllPosts() []Post {
	var posts []Post
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	row, err := db.Query("SELECT * FROM Posts ORDER BY time")
	if err != nil {
		log.Fatalln("Post query failed: ", err.Error())
	}
	defer row.Close()
	for row.Next() {
		var post_id int
		var comment string
		var user int
		var likes int
		var dislikes int
		var time interface{}
		var tags string
		row.Scan(&post_id, &comment, &user, &likes, &dislikes, &time, &tags)
		log.Println("User: ", post_id, " ", comment, " ", user, " ", likes, " ", dislikes, " ", time)
		posts = append(posts, Post{comment, user, likes, dislikes, time})
	}
	return posts
}

func SelectUniqueUser(userName string) User {
	var user User
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	stm, err := db.Prepare("SELECT * FROM Users WHERE name = ?")
	if err != nil {
		log.Fatalln("Statement failed: ", err.Error())
	}
	err = stm.QueryRow(userName).Scan(&user.ID, &user.UUID, &user.Name, &user.Email, &user.Password, &user.Permission)
	if err != nil {
		log.Fatalln("Query Row failed: ", err.Error())
	}
	return user
}
