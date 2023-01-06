package main

import "forum/dbmanagement"

func initaliseDatabase() {
	dbmanagement.CreateDatabase()
	// dbmanagement.InsertUser("0001", "8423479283", "Jupiter", "Koberich-Coles", "387493874")
	// dbmanagement.SelectUser()
}
