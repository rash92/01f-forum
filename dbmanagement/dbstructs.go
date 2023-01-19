package dbmanagement

import "time"

type User struct {
	UUID          string
	Name          string
	Email         string
	Password      string
	Permission    string
	Notifications []Notification
}

type Post struct {
	UUID          string
	Title         string
	Content       string
	OwnerId       string
	Likes         int
	Dislikes      int
	Tag           string
	Time          time.Time
	FormattedTime string
	NumOfComments int
}

type Comment struct {
	UUID          string
	Content       string
	PostId        string
	OwnerId       string
	OwnerName     string
	Likes         int
	Dislikes      int
	Time          time.Time
	FormattedTime string
}

type Session struct {
	UUID      string
	UserId    string
	CreatedAt time.Time
}

type Tag struct {
	UUID    string
	TagName string
}

type AdminRequest struct {
	UUID            string
	RequestFromId   string
	RequestFromName string
	Content         string
}

type Notification struct {
	UUID      string
	Receiver  string
	PostId    string
	CommentId string
	Sender    string
	Reaction  int
	Statement string
}
