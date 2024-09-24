package store

import "time"

type User struct {
	ID    string
	Email string
}

type UserInfo struct {
	UserId    string
	Email     string
	FirstName string
	LastName  string
}

type Post struct {
	ID          int
	UserID      string
	Title       string
	IsAnon      bool
	Description string
	TypeID      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
