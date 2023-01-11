package dbmanagement

import (
	"database/sql"
	"forum/utils"
	"log"
)

func DeleteFromTableWithUUID(table string, UUID string) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	log.Println("Deleting ", table, " record for uuid: ", UUID, "...")

	deleteRowStatement := "DELETE FROM ? WHERE uuid = ?"
	statement, err := db.Prepare(deleteRowStatement)
	utils.HandleError("Delete Prepare failed: ", err)

	_, err = statement.Exec(table, UUID)
	utils.HandleError("Statement Exec failed: ", err)
}
