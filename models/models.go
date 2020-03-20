package models

import (
	"time"
)

const (
	NormalUser   uint = 0 // Suggest books, suggest authors
	ElevatedUser uint = 1 // Validate books, validate authors
	AdminUser    uint = 2 // Admin
)

const (
	ValidateBooks        uint = 0
	ValidateAuthors      uint = 1
	ElevateUser          uint = 2
	DemoteUser           uint = 3
	ValidateSingleAuthor uint = 4
	ValidateSingleBook   uint = 5
)

type User struct {
	Id        uint64
	CreatedAt time.Time
	Email     string
	UserGroup uint64
}

type Book struct {
	Title       string
	Description string
	SuggestedBy *User // Not Null
	ValidatedBy *User
}
