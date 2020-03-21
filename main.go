package main

import (
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
	{
		v1.GET("/clear", web.Clear)
		v1.GET("/users", web.ListUsers)
		v1.GET("/addUser", web.AddUser)
		v1.GET("/addBook", web.AddBook)
		v1.GET("/getUser", web.GetUser)
	}


	router.GET("/", web.ExecuteTemplate("index", "index.html", web.Index))
	router.Run(":80")
}
