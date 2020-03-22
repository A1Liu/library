package web

import (
	"errors"
	"github.com/A1Liu/webserver/database"
	"github.com/A1Liu/webserver/models"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

var (
	AuthorFirstNameInvalid = errors.New("author's first name is invalid")
	AuthorLastNameInvalid = errors.New("author's last name is invalid")
)

func AddAuthorsApi(authors *gin.RouterGroup) {

	authors.GET("/all", func(c *gin.Context) {
		pageIndex, err := strconv.ParseUint(c.Query("pageIndex"), 10, 64)
		if err != nil {
			pageIndex = 0
		}
		authors, err := database.SelectAuthors(pageIndex)
		JsonInfer(c, authors, err)
	})

	authors.GET("/add", func(c *gin.Context) {
		firstName, ok := c.GetQuery("firstName")
		firstName = strings.TrimSpace(firstName)
		if !ok || firstName == "" {
			JsonFail(c, AuthorFirstNameInvalid)
			return
		}

		lastName, ok := c.GetQuery("firstName")
		lastName = strings.TrimSpace(lastName)
		if !ok || lastName == "" {
			JsonFail(c, AuthorLastNameInvalid)
			return
		}

		user, err := QueryParamToken(c)
		if JsonFail(c, err) {
			return
		}

		id, err := database.InsertAuthor(user, firstName, lastName)
		JsonInfer(c, id, err)
	})

	authors.GET("/validate", func(c *gin.Context) {
		user, err := QueryParamToken(c)
		authorId, err := QueryParamUint(c, "authorId")
		if JsonFail(c, err) {
			return
		}

		err = database.ValidateAuthor(user, *authorId)
		JsonInfer(c, nil, err)
	})

	authors.GET("/credit", func(c *gin.Context) {
		user, err := QueryParamToken(c)
		if JsonFail(c, err) {
			return
		}

		authorId, err := QueryParamUint(c, "authorId")
		if JsonFail(c, err) {
			return
		}

		bookId, err := QueryParamUint(c, "bookId")
		if JsonFail(c, err) {
			return
		}

		_, err = database.InsertWrittenBy(user, *authorId, *bookId)
		JsonInfer(c, nil, err)
	})

	authors.GET("/validateCredit", func(c *gin.Context) {
		user, err := QueryParamToken(c)
		if JsonFail(c, err) {
			return
		}

		authorId, err := QueryParamUint(c, "authorId")
		if JsonFail(c, err) {
			return
		}

		bookId, err := QueryParamUint(c, "bookId")
		if JsonFail(c, err) {
			return
		}

		ok, err := database.HasPermissions(user, []models.Permission{
			*models.TargetedPermission(models.ValidateSingleAuthor, *authorId),
			*models.TargetedPermission(models.ValidateSingleBook, *bookId),
		})
		if JsonFail(c, err) {
			return
		}
		if !ok {
			JsonFail(c, MissingPermissions)
			return
		}

		err = database.ValidateWrittenBy(user, *authorId, *bookId)
		JsonInfer(c, nil, err)
	})

	authors.GET("/get", func(c *gin.Context) {
		authorId, err := QueryParamUint(c, "authorId")
		if JsonFail(c, err) {
			return
		}

		author, err := database.GetAuthor(*authorId)
		JsonInfer(c, author, err)
	})

	authors.GET("/merge", func(c *gin.Context) {
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
			*models.TargetedPermission(models.ValidateSingleAuthor, *from),
			*models.TargetedPermission(models.ValidateSingleAuthor, *into),
		})
		if JsonFail(c, err) {
			return
		}
		if !ok {
			JsonFail(c, MissingPermissions)
			return
		}

		err = database.MergeAuthorInto(*from, *into)
		JsonInfer(c, nil, err)
	})
}
