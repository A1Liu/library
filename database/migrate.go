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
const version = 1

var db *sql.DB = nil
var migrater *migrate.Migrate = nil

func getDb() *sql.DB {
	if db != nil {
		return db
	}

	c, err := pgx.ParseURI("psql://webserver:webserver@localhost/webserver")
	if err != nil {
		log.Fatal(err)
	}

	db = stdlib.OpenDB(c)
	return db
}

func GetMigrate() *migrate.Migrate {
	if migrater != nil {
		return migrater
	}

	instance, err := postgres.WithInstance(getDb(), new(postgres.Config))
	if err != nil {
		log.Fatal(err)
	}

	migrater, err = migrate.NewWithDatabaseInstance("file://migrations", "postgres", instance)
	if err != nil {
		log.Fatal(err)
	}

	return migrater
}

func GetDb() *sql.DB {
	if migrater != nil {
		return db
	}

	err := GetMigrate().Migrate(version) // current version
	if err != nil {
		log.Fatal(err)
	}

	return db
}
