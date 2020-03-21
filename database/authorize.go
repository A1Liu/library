package database

import (
	"database/sql"
	"errors"
	"github.com/A1Liu/webserver/models"
	"github.com/A1Liu/webserver/utils"
	sq "github.com/Masterminds/squirrel"
	"log"
	"time"
)

var (
	IncorrectPassword = errors.New("password was incorrect")
	NonexistentUser   = errors.New("user doesn't exist")
	InvalidToken      = errors.New("token is not valid")
	ExpiredToken      = errors.New("token has expired")
)

func AuthorizeWithPassword(db *sql.DB, usernameOrEmail, password string) (*models.User, error) {
	rows, err := psql.Select("*").
		From("users").
		Where(sq.Or{sq.Eq{"username": usernameOrEmail}, sq.Eq{"email": usernameOrEmail}}).
		RunWith(db).
		Query()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, NonexistentUser
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

func AuthorizeWithToken(db *sql.DB, token string) (*models.User, error) {
	if len(token) != 128 {
		return nil, InvalidToken
	}

	rows, err := psql.Select("users.id", "users.created_at", "users.username", "users.email", "users.user_group", "tokens.expires_at").
		From("users").
		Join("tokens ON tokens.user_id = users.id").
		Where(sq.Eq{"tokens.value": token}).
		RunWith(db).
		Query()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, NonexistentUser
	}

	var user models.User
	var tokenExpiry time.Time
	rows.Scan(&user.Id, &user.CreatedAt, &user.Username, &user.Email, &user.UserGroup, &tokenExpiry)
	if rows.Next() {
		log.Fatal("Unique constraint has been broken somehow.")
	}

	if err = rows.Err(); err != nil {
		log.Fatal("Why?", err)
	}

	if time.Now().Before(tokenExpiry) {
		return &user, rows.Err()
	}

	return nil, ExpiredToken
}

func CreateToken(db *sql.DB, userId uint64) (string, error) {
	token := utils.RandomString(128)
	_, err := psql.Insert("tokens").
		Columns("expires_at", "user_id", "value").
		Values(time.Now().Add(time.Hour*24*30), userId, token).
		RunWith(db).
		Exec()

	if err != nil {
		return "", err
	}

	return token, nil
}
