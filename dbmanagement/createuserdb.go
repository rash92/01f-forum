package dbmanagement

import (
	"database/sql"
	"log"
	"os"
)

func CreateUserDb() {
	log.Println("Creating user.db...")
	file, err := os.Create("user.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("user.db create successfully!")

}

func CreatePostsDb() {
	log.Println("Creating posts.db...")
	file, err := os.Create("posts.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("posts.db create successfully!")
}

func CreateTable(db *sql.DB, table string) {

}
