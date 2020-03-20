package main

import (
	"fmt"
	database "github.com/A1Liu/webserver/database"
	_ "github.com/A1Liu/webserver/models"
	"log"
)

func main() {
	database.GetDb()
	err := database.GetMigrate().Down()
	if err != nil {
		log.Fatal(err)
	}
}
