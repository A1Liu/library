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
	MissingPermissions = errors.New("missing permissions")
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
		_, err := database.InsertUser(database.GetDb(), c.Query("username"),
			c.Query("email"), c.Query("email"), 0)
		JsonInfer(c, nil, err)
	})

	users.GET("/token", func(c *gin.Context) {
		user, err := GetQueryParamLogin(c)
		if JsonFail(c, err) {
			return
		}
		token, err := database.CreateToken(database.GetDb(), user.Id)
		JsonInfer(c, token, err)
	})

	users.GET("/get", func(c *gin.Context) {
		user, err := GetQueryParamToken(c)
		isFull, ok := c.GetQuery("full")
		if ok && strings.ToLower(isFull) == "true" {
			// if !JsonFail(c, err) {
			// 	userFull, err := database.GetFullUser(database.GetDb(), user)
			// 	JsonInfer(c, userFull, err)
			// }
			JsonInfer(c, user, err)
		} else {
			JsonInfer(c, user, err)
		}
	})

	users.GET("/permissions/add", func(c *gin.Context) {
		user, err := GetQueryParamToken(c)
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

		if user.UserGroup != models.AdminUser {
			ok, err := database.HasPermissions(database.GetDb(), user,
				[]models.Permission{*models.BroadPermission(models.ElevateUsers), *permission})
			if JsonFail(c, err) {
				return
			}
			if !ok {
				JsonFail(c, MissingPermissions)
				return
			}
		}

		id, err := database.AddPermission(database.GetDb(), user, *target, permission)
		JsonInfer(c, id, err)
	})

	users.GET("/permissions/remove", func(c *gin.Context) {
		user, err := GetQueryParamToken(c)
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

		if user.UserGroup != models.AdminUser {
			ok, err := database.HasPermissions(database.GetDb(), user,
				[]models.Permission{*models.BroadPermission(models.DemoteUsers), *permission})
			if JsonFail(c, err) {
				return
			}
			if !ok {
				JsonFail(c, MissingPermissions)
				return
			}
		}

		err = database.RemovePermissions(database.GetDb(), *target, permission)
		JsonInfer(c, err, err)
	})

}
