package models

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

const (
	NormalUser uint64 = 0 // Suggest books, suggest authors
	AdminUser  uint64 = 1 // Admin
)

var (
	InvalidUserGroup = errors.New("user group was invalid")
)

type User struct {
	Id         uint64
	CreatedAt  time.Time
	Username   string
	Email      string
	UserGroup  uint64
	ProfilePic sql.NullInt64
}

func (user *User) NilId() *uint64 {
	if user == nil {
		return nil
	} else {
		return &user.Id
	}
}

func GetUserGroup(value uint64) string {
	switch value {
	case NormalUser:
		return "NormalUser"
	case AdminUser:
		return "AdminUser"
	default:
		log.Fatal("Got value ", value)
		return ""
	}
}

func IsValidUserGroup(value uint64) bool {
	return value < 3
}

type Book struct {
	Id          uint64
	SuggestedBy *uint64
	SuggestedAt time.Time
	ValidatedBy sql.NullInt64
	ValidatedAt *time.Time
	Title       string
	Description string
	ImageId     sql.NullInt64
}

type Author struct {
	Id          uint64
	SuggestedBy *uint64
	SuggestedAt time.Time
	ValidatedBy sql.NullInt64
	ValidatedAt *time.Time
	FirstName   string
	LastName    string
	ImageId     sql.NullInt64
}

type WrittenBy struct {
	AuthorId uint64
	BookId   uint64
}
