package web

import (
	"github.com/A1Liu/webserver/database"
	"github.com/gin-gonic/gin"
	"strconv"
)

func ListUsers(c *gin.Context) {
	users, err := database.SelectUsers(database.GetDb(), 50, 0)
	jsonInfer(c, users, err)
}

func Clear(c *gin.Context) {
	err := database.Clear()
	jsonInfer(c, err, err)
}

func AddUser(c *gin.Context) {
	userGroup, err := strconv.ParseUint(c.Query("userGroup"), 10, 64)
	username := c.Query("username")
	email := c.Query("email")
	password := c.Query("email")

	token, err := database.InsertUser(database.GetDb(), username, email, password, userGroup)
	jsonInfer(c, token, err)
}

func AddBook(c *gin.Context) {
	title:= c.Query("title")
	description:= c.Query("description")

	user, err := getQueryParamLogin(c)
	if err != nil && err != MissingAuthorization {
		jsonInfer(c, nil, err)
	}
	var id *uint64
	if user == nil {
		id = nil
	} else {
		id = &user.Id
	}

	jsonInfer(c, nil, database.InsertBook(database.GetDb(), id, title, description))
}

func GetUser(c *gin.Context) {
	user, err := getQueryParamLogin(c)
	jsonInfer(c, user, err)
}
