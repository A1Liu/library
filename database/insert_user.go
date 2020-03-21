package database

import (
	"database/sql"
	"errors"
	"github.com/A1Liu/webserver/utils"
	"regexp"
	"strings"
	"time"
)

var (
	InvalidUsername = errors.New("username was invalid")
	InvalidEmail  = errors.New("email was invalid")
	// https://www.alexedwards.net/blog/validation-snippets-for-go#email-validation
	rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

)

func InsertUser(db *sql.DB, username, email, password string, userGroup uint64) (string,error) {
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
	token := utils.RandomString(128)
	_, err := psql.Insert("tokens").
		Columns("expires_at", "user_id", "value").
		Values(time.Now().Add(time.Hour * 24 * 30), id, token).
		RunWith(db).
		Exec()

	if err != nil {
		return "", err
	}

	return token, nil
}
