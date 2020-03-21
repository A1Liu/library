package database

import (
	"database/sql"
	"errors"
	"github.com/A1Liu/webserver/models"
	sq "github.com/Masterminds/squirrel"
	"regexp"
	"strings"
)

var (
	InvalidUsername = errors.New("username was invalid")
	InvalidEmail    = errors.New("email was invalid")
	// https://www.alexedwards.net/blog/validation-snippets-for-go#email-validation
	InvalidPageSize = errors.New(
		"gave an invalid page size. Pages can between 50 and 100 entries long")

	rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func InsertUser(db *sql.DB, username, email, password string, userGroup uint64) (string, error) {
	if len(username) > 16 || strings.Contains(username, "@") {
		return "", InvalidUsername
	} else if len(email) > 254 || !rxEmail.MatchString(email) {
		return "", InvalidEmail
	}

	row := psql.Insert("users").
		Columns("username", "email", "password", "user_group").
		Values(username, email, password, userGroup).
		Suffix("RETURNING \"id\"").
		RunWith(db).
		QueryRow()

	var id uint64
	row.Scan(&id)
	return CreateToken(db, id)
}

func SelectUsers(db *sql.DB, pageIndex uint64) ([]models.User, error) {

	users := make([]models.User, 50)[:0]

	rows, err := psql.Select("id", "created_at", "username", "email", "user_group").
		From("users").
		Where(sq.Lt{"id": 50 * (pageIndex + 1)}).
		Where(sq.GtOrEq{"id": 50 * pageIndex}).
		Limit(50).
		RunWith(db).
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
