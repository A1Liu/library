package database

import (
	"database/sql"
)

func InsertBook(db *sql.DB, userId *uint64, title string, description string) error {
	rows, err := psql.Insert("books").
				Columns("suggested_by", "title", "description").
				Values(userId, title, description).
				RunWith(db).
				Query()
	if err != nil {
		return err
	}
	defer rows.Close()
	return rows.Err()
}


