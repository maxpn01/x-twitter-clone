package models

import "time"

type User struct {
	ID           string
	Email        string
	Username     string
	Fullname     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Post struct {
	ID        string
	Text      string
	CreatedAt time.Time
	UserID    string
	Username  string
	Fullname  string
}

type Follows struct {
	FollowerID  string
	FollowingID string
	CreatedAt   time.Time
}
