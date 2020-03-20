package database

import (
	"database/sql"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"log"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// version defines the current migration version. This ensures the app
// is always compatible with the version of the database.
const COMPAT_VERSION = 1

var db *sql.DB = nil
var migrater *migrate.Migrate = nil

func getDb() (*sql.DB, *migrate.Migrate) {
	if migrater != nil {
		return db, migrater
	}

	c, err := pgx.ParseURI("psql://webserver:webserver@localhost/webserver")
	if err != nil {
		log.Fatal(err)
	}

	db = stdlib.OpenDB(c)

	instance, err := postgres.WithInstance(db, new(postgres.Config))
	if err != nil {
		log.Fatal(err)
	}

	migrater, err = migrate.NewWithDatabaseInstance("file://migrations", "postgres", instance)
	if err != nil {
		log.Fatal(err)
	}

	return db, migrater
}

func ExperimentDb(try func(*sql.DB)) {
	db, migrater := getDb()
	originalVersion, dirty, err := migrater.Version()
	if dirty {
		log.Fatal("Database is dirty with version ", originalVersion)
	}
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatal(err)
	}

	err = migrater.Migrate(COMPAT_VERSION)
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	try(db)

	migrater.Migrate(originalVersion)
}
