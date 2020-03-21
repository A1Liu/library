package main

import (
	"github.com/A1Liu/webserver/database"
	"github.com/gin-gonic/gin"
	"log"
	// "github.com/A1Liu/webserver/models"
	"github.com/A1Liu/webserver/web"
	sq "github.com/Masterminds/squirrel"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	router := gin.Default()

	v1 := router.Group("/api/v1")
	users := v1.Group("/users")
	web.AddUsersApi(users)

	permissions := users.Group("/permissions")
	web.AddPermissionsApi(permissions)

	books := v1.Group("/books")
	web.AddBooksApi(books)

	v1.GET("/clear", func(c *gin.Context) {
		err := database.Clear()
		web.JsonInfer(c, err, err)
	})

	router.GET("/", web.ExecuteTemplate("index", "index.html", web.Index))
	router.Run(":80")
}
