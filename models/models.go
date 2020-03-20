package models

const (
  NormalUser int = 0 // Suggest books
  ElevatedUser int = 1 // Validate books
  AdminUser int = 2 // Admin
)

type User struct {
  email string
  user_group int
}
