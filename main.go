package main

import (
	"fmt"
	database "github.com/A1Liu/webserver/database"
	"log"
	"net/http"

	// "github.com/A1Liu/webserver/models"
	"github.com/A1Liu/webserver/web"
	sq "github.com/Masterminds/squirrel"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func main() {

	fmt.Println("Inserting into database...")
	_, err := psql.Insert("").
		Into("users").
		Columns("email", "password", "user_group").
		Values("hi@gmil.com", "hellofresh", 1).
		RunWith(database.GetDb()).
		Query()

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api/v1/clear", web.Clear)
	http.HandleFunc("/api/v1/users", web.ListUsers)
	http.ListenAndServe(":8080", nil)
}
