package database

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"
)

var (
	InvalidUsername = errors.New("username was invalid")
	InvalidEmail  = errors.New("email was invalid")
	// https://www.alexedwards.net/blog/validation-snippets-for-go#email-validation
	rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

)



func InsertUser(db *sql.DB, username, email, password string, userGroup uint64) error {
	if strings.Contains(username, "@") {
		return InvalidUsername
	} else if len(email) > 254 || !rxEmail.MatchString(email) {
		return InvalidEmail
	}

	rows, err := psql.Insert("users").
		Columns("username", "email", "password", "user_group").
		Values(username, email, password, userGroup).
		RunWith(db).
		Query()
	if err == nil {
		rows.Close()
	}
	return err
}
