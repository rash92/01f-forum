package main

import "forum/dbmanagement"

func main() {
	dbmanagement.CreateUserDb()
	dbmanagement.CreatePostsDb()
}
