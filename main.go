package main

import (
	// "log"
	database "github.com/A1Liu/webserver/database"
)

func main() {
	database.GetDb()
	// err := database.GetMigrate().Down()
	// if err != nil {
	//   log.Fatal(err)
	// }
}
