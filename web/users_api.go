package web

import (
	"errors"
	"github.com/A1Liu/webserver/database"
	"github.com/A1Liu/webserver/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

var (
	MissingPermissions = errors.New("missing permissions")
)

func AddUsersApi(users *gin.RouterGroup) {
	users.GET("/all", func(c *gin.Context) {
		pageIndex, err := strconv.ParseUint(c.Query("pageIndex"), 10, 64)
		if err != nil {
			pageIndex = 0
		}
		users, err := database.SelectUsers(pageIndex)
		JsonInfer(c, users, err)
	})

	users.GET("/add", func(c *gin.Context) {
		_, err := database.InsertUser(c.Query("username"),
			c.Query("email"), c.Query("password"), 0)
		JsonInfer(c, nil, err)
	})

	users.GET("/token", func(c *gin.Context) {
		user, err := QueryParamLogin(c)
		if JsonFail(c, err) {
			return
		}
		token, err := database.CreateToken(user.Id)
		JsonInfer(c, token, err)
	})

	users.GET("/get", func(c *gin.Context) {
		user, err := QueryParamToken(c)
		JsonInfer(c, user, err)
	})
}

func AddPermissionsApi(permissions *gin.RouterGroup) {
	permissions.GET("/add", func(c *gin.Context) {
		user, err := QueryParamToken(c)
		if JsonFail(c, err) {
			return
		}

		target, err := QueryParamUint(c, "target")
		if JsonFail(c, err) {
			return
		}

		reference, err := QueryParamUint(c, "reference")
		if JsonFail(c, err) {
			return
		}

		permission, err := models.BuildPermission(c.Query("permission"), *reference)
		if err != nil {
			JsonFail(c, err)
			return
		}

		ok, err := database.HasPermissions(user,
			[]models.Permission{*models.BroadPermission(models.ElevateUsers), *permission})
		if JsonFail(c, err) {
			return
		}
		if !ok {
			JsonFail(c, MissingPermissions)
			return
		}

		id, err := database.AddPermission(user, *target, permission)
		JsonInfer(c, id, err)
	})

	permissions.GET("/remove", func(c *gin.Context) {
		user, err := QueryParamToken(c)
		if JsonFail(c, err) {
			return
		}

		target, err := QueryParamUint(c, "target")
		if JsonFail(c, err) {
			return
		}

		reference, err := QueryParamUint(c, "reference")
		if JsonFail(c, err) {
			return
		}

		permission, err := models.BuildPermission(c.Query("permission"), *reference)
		if err != nil {
			JsonFail(c, err)
			return
		}

		ok, err := database.HasPermissions(user,
			[]models.Permission{*models.BroadPermission(models.DemoteUsers), *permission})
		if JsonFail(c, err) {
			return
		}
		if !ok {
			JsonFail(c, MissingPermissions)
			return
		}

		err = database.RemovePermissions(*target, *permission)
		JsonInfer(c, err, err)
	})
}
