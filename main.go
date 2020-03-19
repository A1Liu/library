package main

import (
	database "github.com/A1Liu/webserver/database"
)

func main() {
	database.GetMigrate().Down()
}
