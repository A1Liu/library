package web

import (
	"github.com/A1Liu/webserver/database"
	"github.com/gin-gonic/gin"
	"strconv"
)

func AddUsersApi(users *gin.RouterGroup) {
	users.GET("/all", func(c *gin.Context) {
		pageIndex, err := strconv.ParseUint(c.Query("pageIndex"), 10, 64)
		if err != nil {
			pageIndex = 0
		}
		users, err := database.SelectUsers(database.GetDb(), pageIndex)
		JsonInfer(c, users, err)
	})

	users.GET("/add", func(c *gin.Context) {
		token, err := database.InsertUser(database.GetDb(), c.Query("username"),
			c.Query("email"), c.Query("email"), 0)
		JsonInfer(c, token, err)
	})

	users.GET("/token", func(c *gin.Context) {
		user, err := GetQueryParamLogin(c)
		if err != nil {
			JsonFail(c, err)
			return
		}
		token, err := database.CreateToken(database.GetDb(), user.Id)
		JsonInfer(c, token, err)
	})

	users.GET("/get", func(c *gin.Context) {
		user, err := GetQueryParamToken(c)
		JsonInfer(c, user, err)
	})
}
