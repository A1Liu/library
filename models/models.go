package models

const (
	NormalUser   int = 0 // Suggest books
	ElevatedUser int = 1 // Validate books
	AdminUser    int = 2 // Admin
)

type User struct {
  Id        uint
	CreatedAt string
	Email     string
	UserGroup uint
}

type Book struct {
	Title       string
	Description string
	SuggestedBy *User // Not Null
	ValidatedBy *User
}
