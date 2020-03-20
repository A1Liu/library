package main

import (
	"database/sql"
	"fmt"
	database "github.com/A1Liu/webserver/database"
	// "github.com/A1Liu/webserver/models"
	sq "github.com/Masterminds/squirrel"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func main() {

	database.ExperimentDb(func(db *sql.DB) error {

		fmt.Println("Inserting into database...")
		_, err := psql.Insert("").
			Into("users").
			Columns("email", "password", "user_group").
			Values("hi@gmail.com", "hellofresh", 0).
			RunWith(db).
			Query()

		if err != nil {
			return err
		}

		fmt.Println("Reading inserted data from database...")
		users, err := database.SelectUsers(db, 50, 0)
		if err != nil {
			return err
		}

		fmt.Println(users)

		return nil
	})
}
