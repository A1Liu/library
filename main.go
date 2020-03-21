package main

import (
	"log"
	"net/http"

	// "github.com/A1Liu/webserver/models"
	"github.com/A1Liu/webserver/web"
	sq "github.com/Masterminds/squirrel"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	http.HandleFunc("/api/v1/clear", web.Clear)
	http.HandleFunc("/api/v1/users", web.ListUsers)
	http.HandleFunc("/api/v1/addUser", web.AddUser)
	http.HandleFunc("/api/v1/addBook", web.AddBook)
	http.HandleFunc("/", web.ExecuteTemplate("index", "index.html", web.Index))
	http.ListenAndServe(":80", nil)
}
