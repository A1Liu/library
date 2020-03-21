package web

import (
	"github.com/A1Liu/webserver/database"
	"github.com/A1Liu/webserver/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

func AddBooksApi(books *gin.RouterGroup) {
	books.GET("/all", func(c *gin.Context) {
		pageIndex, err := strconv.ParseUint(c.Query("pageIndex"), 10, 64)
		if err != nil {
			pageIndex = 0
		}
		books, err := database.SelectBooks(database.GetDb(), pageIndex)
		JsonInfer(c, books, err)
	})

	books.GET("/add", func(c *gin.Context) {
		title := c.Query("title")
		description := c.Query("description")

		user, err := GetQueryParamLogin(c)
		if err != nil && err != NoLoginInformation {
			JsonFail(c, err)
			return
		}

		if user.UserGroup == models.AdminUser {

		}

		bookId, err := database.InsertBook(database.GetDb(), user, title, description)
		JsonInfer(c, bookId, err)
	})

}
