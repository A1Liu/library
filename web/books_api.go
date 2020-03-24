package web

import (
	"github.com/A1Liu/library/database"
	"github.com/A1Liu/library/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

func AddBooksApi(books *gin.RouterGroup) {
	books.GET("/all", func(c *gin.Context) {
		pageIndex, err := strconv.ParseUint(c.Query("pageIndex"), 10, 64)
		if err != nil {
			pageIndex = 0
		}
		books, err := database.SelectBooks(pageIndex)
		JsonInfer(c, books, err)
	})

	books.POST("/add", func(c *gin.Context) {
		title := c.Query("title")
		description := c.Query("description")

		user, err := QueryParamToken(c)
		if err == NoLoginInformation {
			user = nil
		} else if JsonFail(c, err) {
			return
		}
		bookId, err := database.InsertBook(user, title, description)
		JsonInfer(c, bookId, err)
	})

	books.POST("/validate", func(c *gin.Context) {
		user, err := QueryParamToken(c)
		bookId, err := QueryParamUint(c, "bookId")
		if JsonFail(c, err) {
			return
		}

		err = database.ValidateAuthor(user, *bookId)
		JsonInfer(c, nil, err)
	})

	books.GET("/get", func(c *gin.Context) {
		bookId, err := QueryParamUint(c, "bookId")
		if JsonFail(c, err) {
			return
		}

		book, err := database.GetBook(*bookId)
		JsonInfer(c, book, err)
	})

	books.POST("/merge", func(c *gin.Context) {
		from, err := QueryParamUint(c, "from")
		if JsonFail(c, err) {
			return
		}

		into, err := QueryParamUint(c, "into")
		if JsonFail(c, err) {
			return
		}

		user, err := QueryParamToken(c)
		if JsonFail(c, err) {
			return
		}

		ok, err := database.HasPermissions(user, []models.Permission{
			*models.TargetedPermission(models.ValidateSingleBook, *from),
			*models.TargetedPermission(models.ValidateSingleBook, *into),
		})
		if JsonFail(c, err) {
			return
		}
		if !ok {
			JsonFail(c, MissingPermissions)
			return
		}

		err = database.MergeBookInto(*from, *into)
		JsonInfer(c, nil, err)
	})

}
