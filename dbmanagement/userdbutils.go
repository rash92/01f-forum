package dbmanagement

import (
	"database/sql"
	"forum/utils"
	"log"
)

/*
Generates a new user in the database.  The UUID is generated internally here and stored to the database (this can also be referred to as the userID).

The inserted User is also returned in case it is needed to be used straight away but it is not necessary.
*/
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

func UpdateUserPermissionFromUUID(UUID string, newpermission string) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	log.Println("updating user permission to: ", newpermission)

	updateUserData := "UPDATE Users SET permission = ? WHERE uuid = ?"
	statement, err := db.Prepare(updateUserData)
	utils.HandleError("User Update Prepare failed: ", err)

	_, err = statement.Exec(newpermission, UUID)
	utils.HandleError("Statement Exec failed: ", err)
}

func UpdateUserPermissionFromName(Name string, newpermission string) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	log.Println("updating user permission to: ", newpermission)

	updateUserData := "UPDATE Users SET permission = ? WHERE name = ?"
	statement, err := db.Prepare(updateUserData)
	utils.HandleError("User Update Prepare failed: ", err)

	_, err = statement.Exec(newpermission, Name)
	utils.HandleError("Statement Exec failed: ", err)
}

/*
Used to display all currently registered users.  Should only be used internally as information is not relevant for the website.
*/
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

func SelectAllUsers() []User {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	row, err := db.Query("SELECT * FROM Users")
	utils.HandleError("User query failed: ", err)
	defer row.Close()

	var allUsers []User
	for row.Next() {
		var currentUser User
		row.Scan(&currentUser.UUID, &currentUser.Name, &currentUser.Email, &currentUser.Password, &currentUser.Permission)
		allUsers = append(allUsers, currentUser)
	}
	return allUsers
}

/*
Initially used for when a user is trying to log in.  Returns a User's information when searched for by name.
*/
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

/*
Could be used for if a user wanted to log in using their email address.  Returns a User's information when searched for by email.
*/
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

/*
Used when you have the users UUID (userID).  For example, within a session (displaying user information such as username), or when displaying post and comment details.
*/
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

/*
Gets the user using the current session.  Used to assign the correct userID if a user posts, likes, dislikes, or comments.
*/
func SelectUserFromSession(UUID string) User {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	var userID string
	err := db.QueryRow("SELECT userId FROM Sessions WHERE uuid = ?", UUID).Scan(&userID)
	utils.HandleError("User from session query failed: ", err)

	var user User
	err = db.QueryRow("SELECT * FROM Users WHERE uuid = ?", userID).Scan(&user.UUID, &user.Name, &user.Email, &user.Password, &user.Permission)
	utils.HandleError("User query failed: ", err)

	return user
}
