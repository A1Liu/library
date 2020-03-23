package database

import (
	"errors"
	"github.com/A1Liu/webserver/models"
	sq "github.com/Masterminds/squirrel"
	"log"
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

func InsertUser(username, email, password string, userGroup uint64) (uint64, error) {
	log.Println("new user", username, email)
	username = strings.ToLower(username)
	if len(username) < 2 || len(username) > 16 || !rxUsername.MatchString(username) {
		return 0, InvalidUsername
	} else if len(email) > 254 || !rxEmail.MatchString(email) {
		return 0, InvalidEmail
	} else if !models.IsValidUserGroup(userGroup) {
		return 0, models.InvalidUserGroup
	} else if len(password) < 5 || len(password) > 32 {
		return 0, InvalidPassword
	}

	row := psql.Insert("users").
		Columns("username", "email", "password", "user_group").
		Values(username, email, password, userGroup).
		Suffix("RETURNING \"id\"").
		RunWith(globalDb).
		QueryRow()

	var id uint64
	err := row.Scan(&id)
	return id, err
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
