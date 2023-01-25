package dbmanagement

import (
	"database/sql"
	"forum/utils"
)

func InsertTag(tag string) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	utils.WriteMessageToLogFile("Inserting tag record...")

	UUID := GenerateUUIDString()
	insertPostData := "INSERT INTO Tags(uuid, tagname) VALUES (?, ?)"
	statement, err := db.Prepare(insertPostData)
	utils.HandleError("User Prepare failed", err)

	_, err = statement.Exec(UUID, tag)
	utils.HandleError("Statement Exec failed", err)
}

func SelectAllTags() []Tag {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	row, err := db.Query("SELECT * FROM Tags")
	utils.HandleError("All tags query failed", err)
	defer row.Close()

	var allTags []Tag
	for row.Next() {
		var currentTag Tag
		row.Scan(&currentTag.UUID, &currentTag.TagName)
		allTags = append(allTags, currentTag)
	}
	return allTags
}
