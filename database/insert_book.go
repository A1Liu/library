package database

import (
	"database/sql"
	"github.com/A1Liu/webserver/models"
)

func InsertBook(db *sql.DB, user models.User, title string, description string) error {
	rows, err := psql.Insert("books").
				Columns("suggested_by", "title", "description").
				Values(user.Id, title, description).
				RunWith(db).
				Query()
	if err != nil {
		rows.Close()
	}
	return err
}


