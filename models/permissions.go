package models

const (
	ValidateBooks        uint = 0
	ValidateAuthors      uint = 1
	ElevateUser          uint = 2
	DemoteUser           uint = 3
	ValidateSingleAuthor uint = 4
	ValidateSingleBook   uint = 5
)

type Permission struct {
	PermissionType uint64
	Reference      uint64
}
