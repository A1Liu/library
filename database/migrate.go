package users

import (
	"database/sql"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"log"

	sq "github.com/Masterminds/squirrel"
	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// version defines the current migration version. This ensures the app
// is always compatible with the version of the database.
const version = 1

var db *sql.DB = func() *sql.DB {
	c, err := pgx.ParseURI("psql://webserver:webserver@localhost/webserver")
	if err != nil {
		log.Fatal(err)
	}

	db := stdlib.OpenDB(c)
	targetInstance, err := postgres.WithInstance(db, new(postgres.Config))
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", targetInstance)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Migrate(version) // current version
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	return db
}()

func Insert(into string) sq.InsertBuilder {
	return sq.Insert(into).RunWith(db)
}
