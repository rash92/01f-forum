package main

import "forum/dbmanagement"

func main() {
	dbmanagement.CreateDatabase()
	dbmanagement.InsertUser(0001, "8423479283", "Jupiter", "Koberich-Coles", "387493874", "mod")
	dbmanagement.SelectUser()
}
