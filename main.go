package main

import (
	"database/sql"
	"fmt"
	database "github.com/A1Liu/webserver/database"
	"github.com/A1Liu/webserver/models"
	sq "github.com/Masterminds/squirrel"
	"log"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func main() {
	database.ExperimentDb(func(db *sql.DB) {
		fmt.Println("Inserting into database...")
		builder := psql.Insert("").
			Into("users").
			Columns("email", "password", "user_group").
			Values("hi@gmail.com", "hellofresh", 0)

		fmt.Println(builder.ToSql())

		_, err := builder.
			RunWith(db).
			Query()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Reading inserted data from database...")
		rows, err := psql.Select("email", "user_group").From("users").RunWith(db).Query()
		defer rows.Close()
		if err != nil {
			log.Fatal(err)
		}

		var user models.User

		for rows.Next() {
			err := rows.Scan(&user.Email, &user.UserGroup)

			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Println(user)
			}

		}
	})
}
