package database

import "database/sql"

func InsertUser(db *sql.DB, email string, password string, userGroup uint64) error {
	psql.Insert("").
		Into("users").
		Columns("email", "password", "user_group").
		Values(email, password, userGroup).
		RunWith(db).
		Query()
}
