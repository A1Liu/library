package database

import (
	"database/sql"
	"errors"
	"github.com/A1Liu/webserver/models"
	sq "github.com/Masterminds/squirrel"
	"log"
)

var (
	IncorrectPassword = errors.New("password was incorrect")
	NonexistentUser  = errors.New("user doesn't exist")
	InvalidToken = errors.New("token is not valid")
	ExpiredToken = errors.New("token has expired")
)

func AuthorizeWithPassword(db *sql.DB, usernameOrEmail, password string) (*models.User, error) {
	rows, err := psql.Select("*").
		From("users").
		Where(sq.Or{sq.Eq{"username": usernameOrEmail},sq.Eq{"email": usernameOrEmail}}).
		RunWith(db).
		Query()

	if err != nil {
		return nil,err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil,NonexistentUser
	}

	var user models.User
	var correctPassword string
	rows.Scan(&user.Id, &user.CreatedAt, &user.Username, &user.Email, &correctPassword, &user.UserGroup)
	if rows.Next() {
		log.Fatal("Unique constraint has been broken somehow.")
	}

	if err = rows.Err(); err != nil {
		log.Fatal("Why?", err)
	}

	if correctPassword == password {
		return &user, rows.Err()
	}

	return nil, IncorrectPassword
}

func AuthorizeWithToken(token string) (*models.User, error) {
	log.Fatal("Not implmented")
	return nil, nil
}
