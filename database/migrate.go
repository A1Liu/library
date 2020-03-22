package database

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"log"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// version defines the current migration version. This ensures the app
// is always compatible with the version of the database.
const CompatVersion = 1

var migrated = false

var (
	globalDb      *sql.DB          = nil
	globalMigrate *migrate.Migrate = nil
	psql                           = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

func getMigrateInstance(db *sql.DB) *migrate.Migrate {
	driver, err := postgres.WithInstance(db, new(postgres.Config))
	if err != nil {
		log.Fatal(err)
	}

	migrateInstance, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	return migrateInstance
}

func getDb() (*sql.DB, *migrate.Migrate) {
	if globalDb != nil {
		return globalDb, globalMigrate
	}

	connStr := "postgres://webserver:webserver@localhost/webserver?sslmode=disable"
	var err error
	globalDb, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	globalMigrate = getMigrateInstance(globalDb)
	return globalDb, globalMigrate
}

func ExperimentDb(try func(*sql.DB) error) {
	db, migrater := getDb()
	var originalVersion *uint = nil
	originalVersionValue, dirty, err := migrater.Version()
	if dirty {
		log.Println("Database is dirty with version ", originalVersion)
		migrater.Drop()
	} else if err == migrate.ErrNilVersion {
		originalVersion = nil
	} else if err != nil {
		log.Fatal("Error when getting version of database: ", err)
	} else {
		originalVersion = &originalVersionValue
	}

	log.Println("Migrating...")
	err = migrater.Migrate(CompatVersion)
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No migration performed; actions may be permanent.")
		} else {
			log.Fatal("Error when migrating: ", err)
		}
	}

	err = try(db)
	if err != nil {
		log.Println("Error while executing closure: ", err)
	}

	var err2 error
	if originalVersion == nil {
		err2 = migrater.Down()
	} else {
		err2 = migrater.Migrate(*originalVersion)
	}

	if err2 != nil && err2 != migrate.ErrNoChange {
		log.Println("Error migrating back to version ", originalVersion, ": ", err2)
	}
}

func CommitDbMigrate(try func(*sql.DB) error) {
	db, migrater := getDb()
	var originalVersion *uint = new(uint)
	didMigration := false
	if !migrated {
		originalVersionInner, dirty, err := migrater.Version()
		*originalVersion = originalVersionInner
		if dirty {
			log.Fatal("Database is dirty with version ", originalVersion)
		}
		if err != nil && err != migrate.ErrNilVersion {
			if err == migrate.ErrNilVersion {
				originalVersion = nil
			}
			log.Fatal(err)
		}

		err = migrater.Migrate(CompatVersion)
		if err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		} else if err != migrate.ErrNoChange {
			didMigration = true
		}
	}

	err := try(db)
	if err != nil {
		var err2 error = nil
		if didMigration {
			if originalVersion == nil {
				err2 = migrater.Down()
			} else {
				err2 = migrater.Migrate(*originalVersion)
			}
		}

		if err2 != nil {
			log.Fatal("Error: ", err2, " while processing error: ", err)
		} else {
			log.Fatal(err)
		}
	}
}

func Clear() error {
	if globalDb == nil {
		getDb()
	}

	err := globalMigrate.Drop()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	globalMigrate = getMigrateInstance(globalDb)
	err = globalMigrate.Migrate(CompatVersion)
	if err != nil {
		return err
	}
	return nil
}

func ConnectToDb() {
	if globalDb != nil {
		return
	}

	getDb()
	err := globalMigrate.Migrate(CompatVersion)
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("Error when migrating: ", err)
	}
}
