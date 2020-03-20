package main

import (
	"net/http"

	// "github.com/A1Liu/webserver/models"
	"github.com/A1Liu/webserver/web"
	sq "github.com/Masterminds/squirrel"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func main() {
	http.HandleFunc("/api/v1/clear", web.Clear)
	http.HandleFunc("/api/v1/users", web.ListUsers)
	http.HandleFunc("/api/v1/addUser", web.AddUser)
	http.ListenAndServe(":8080", nil)
}
