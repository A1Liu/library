package database

import (
	"database/sql"
	"errors"
	models "github.com/A1Liu/webserver/models"
	sq "github.com/Masterminds/squirrel"
)

var (
	InvalidPageSize = errors.New(
		"gave an invalid page size. Pages can between 50 and 100 entries long")
)

func SelectUsers(db *sql.DB, pageSize, pageIndex uint64) ([]models.User, error) {
	if pageSize < 50 || pageSize > 250 {
		return nil, InvalidPageSize
	}

	users := make([]models.User, pageSize)[:0]

	rows, err := psql.Select("id", "created_at", "email", "user_group").From("users").
		Where(sq.Lt{"id": pageSize * (pageIndex + 1)}).
		Where(sq.GtOrEq{"id": pageSize * pageIndex}).
		Limit(pageSize).
		RunWith(db).
		Query()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user models.User

	for rows.Next() {
		err := rows.Scan(&user.Id, &user.CreatedAt, &user.Email, &user.UserGroup)
		if err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}
