package database

import (
	"errors"
	"github.com/A1Liu/webserver/models"
	"github.com/A1Liu/webserver/utils"
	sq "github.com/Masterminds/squirrel"
	"time"
)

var (
	IncorrectPassword = errors.New("password was incorrect")
	NonexistentUser   = errors.New("user doesn't exist")
	InvalidToken      = errors.New("token is not valid")
	ExpiredToken      = errors.New("token has expired")
)

func AuthorizeWithPassword(usernameOrEmail, password string) (*models.User, error) {
	row := psql.Select("*").
		From("users").
		Where(sq.Or{sq.Eq{"username": usernameOrEmail}, sq.Eq{"email": usernameOrEmail}}).
		RunWith(globalDb).
		QueryRow()

	var user models.User
	var correctPassword string
	err := row.Scan(&user.Id, &user.CreatedAt, &user.Username,
		&user.Email, &correctPassword, &user.UserGroup)

	if err != nil {
		return nil, err
	}

	// @TODO Change this to use password hashing, salting, etc.
	if correctPassword == password {
		return &user, nil
	}

	return nil, IncorrectPassword
}

func AuthorizeWithToken(token string) (*models.User, error) {
	if len(token) != 128 {
		return nil, InvalidToken
	}

	row := psql.Select("users.id", "users.created_at", "users.username", "users.email", "users.user_group", "tokens.expires_at").
		From("users").
		Join("tokens ON tokens.user_id = users.id").
		Where(sq.Eq{"tokens.value": token}).
		RunWith(globalDb).
		QueryRow()

	var user models.User
	var tokenExpiry time.Time
	err := row.Scan(&user.Id, &user.CreatedAt, &user.Username, &user.Email,
		&user.UserGroup, &tokenExpiry)
	if err != nil {
		return nil, err
	}

	if time.Now().Before(tokenExpiry) {
		return &user, nil
	}

	return nil, ExpiredToken
}

func CreateToken(userId uint64) (string, error) {
	token := utils.RandomString(128)
	_, err := psql.Insert("tokens").
		Columns("expires_at", "user_id", "value").
		Values(time.Now().Add(time.Hour*24*30), userId, token).
		RunWith(globalDb).
		Exec()

	if err != nil {
		return "", err
	}

	return token, nil
}
