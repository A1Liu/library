package database

import (
	"errors"
	"github.com/A1Liu/library/models"
	sq "github.com/Masterminds/squirrel"
	"regexp"
	"strings"
)

var (
	InvalidUsername = errors.New("username was invalid; must be 2-16 long, begin with a letter, and contain alphanumerics/underscores only")
	InvalidEmail    = errors.New("email was in an invalid format")
	InvalidPassword = errors.New("password was invalid; must be 5-32 long")
	// https://www.alexedwards.net/blog/validation-snippets-for-go#email-validation
	InvalidPageSize = errors.New(
		"gave an invalid page size. Pages can between 50 and 100 entries long")

	rxUsername = regexp.MustCompile("^[a-z][a-z0-9_]*$")
	rxEmail    = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func validateUser(username, email string) error {
	if len(username) < 2 || len(username) > 16 || !rxUsername.MatchString(username) {
		return InvalidUsername
	} else if len(email) > 254 || !rxEmail.MatchString(email) {
		return InvalidEmail
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 5 || len(password) > 32 {
		return InvalidPassword
	}
	return nil
}

func UpdateUser(id uint64, username, email string) error {
	username = strings.ToLower(username)
	err := validateUser(username, email)
	if err != nil {
		return err
	}

	_, err = psql.Update("users").
		Set("username", username).
		Set("email", email).
		Where(sq.Eq{"id": id}).
		RunWith(globalDb).
		Exec()
	return err
}

func UpdateUserGroup(id uint64, userGroup uint64) error {
	if !models.IsValidUserGroup(userGroup) {
		return models.InvalidUserGroup
	}

	_, err := psql.Update("users").
		Set("user_group", userGroup).
		Where(sq.Eq{"id": id}).
		RunWith(globalDb).
		Exec()

	return err
}

func UpdateUserPassword(user *models.User, oldPass, newPass string) error {
	err := validatePassword(newPass)
	if err != nil {
		return err
	}

	row := psql.Select("password").
		From("users").
		Where(sq.Eq{"id": user.Id}).
		RunWith(globalDb).
		QueryRow()

	var password string
	err = row.Scan(&password)
	if err != nil {
		return err
	}

	if password != oldPass {
		return IncorrectPassword
	}

	_, err = psql.Update("users").
		Set("password", newPass).
		Where(sq.Eq{"id": user.Id}).
		RunWith(globalDb).
		Exec()
	return err
}

func InsertUser(username, email, password string, userGroup uint64) error {
	username = strings.ToLower(username)
	err := validateUser(username, email)
	if err != nil {
		return err
	}

	err = validatePassword(password)
	if err != nil {
		return err
	}

	if !models.IsValidUserGroup(userGroup) {
		return models.InvalidUserGroup
	}

	_, err = psql.Insert("users").
		Columns("username", "email", "password", "user_group").
		Values(username, email, password, userGroup).
		RunWith(globalDb).
		Exec()
	return err
}

func SelectUsers(pageIndex uint64) ([]models.User, error) {

	users := make([]models.User, 50)[:0]

	rows, err := psql.Select("id", "created_at", "username", "email", "user_group").
		From("users").
		Where(sq.Lt{"id": 50 * (pageIndex + 1)}).
		Where(sq.GtOrEq{"id": 50 * pageIndex}).
		Limit(50).
		RunWith(globalDb).
		Query()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user models.User
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.CreatedAt, &user.Username, &user.Email, &user.UserGroup)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func GetUser(id uint64) (*models.User, error) {
	row := psql.Select("id", "created_at", "username", "email", "user_group").
		From("users").
		Where(sq.Eq{"id": id}).
		RunWith(globalDb).
		QueryRow()

	var user models.User
	err := row.Scan(&user.Id, &user.CreatedAt, &user.Username, &user.Email, &user.UserGroup)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
