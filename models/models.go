package models

const (
	NormalUser   int = 0 // Suggest books
	ElevatedUser int = 1 // Validate books
	AdminUser    int = 2 // Admin
)

type User struct {
	Email     string
	Password  string
	UserGroup int
}

type Book struct {
	Title       string
	Description string
	SuggestedBy *User // Not Null
	ValidatedBy *User
}
