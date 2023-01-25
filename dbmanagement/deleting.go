package dbmanagement

import (
	"database/sql"
	"forum/utils"
)

func DeleteFromTableWithUUID(table string, UUID string) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	utils.WriteMessageToLogFile("Deleting " + table + " record for uuid: " + UUID + "...")

	deleteRowStatement := "DELETE FROM " + table + " WHERE uuid = ?"
	statement, err := db.Prepare(deleteRowStatement)
	utils.HandleError("Delete Prepare failed: ", err)

	_, err = statement.Exec(UUID)
	utils.HandleError("Statement Exec failed: ", err)
}
