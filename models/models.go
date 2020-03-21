package models

import (
	"time"
)

const (
	NormalUser   uint = 0 // Suggest books, suggest authors
	ElevatedUser uint = 1 // Validate books, validate authors
	AdminUser    uint = 2 // Admin
)

type User struct {
	Id        uint64
	CreatedAt time.Time
	Username  string
	Email     string
	UserGroup uint64
}

type Book struct {
	Id          uint64
	Title       string
	Description string
	SuggestedBy *uint64
	SuggestedAt time.Time
	ValidatedBy *uint64
	ValidatedAt *time.Time
}
