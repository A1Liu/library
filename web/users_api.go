package web

import (
	"errors"
	"github.com/A1Liu/library/database"
	"github.com/A1Liu/library/models"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
)

var (
	MissingPermissions = errors.New("missing permissions")
	MissingEmail       = errors.New("missing email")
	MissingUsername    = errors.New("missing username")
	EmptyUpdate        = errors.New("updating nothing")
	GivingSelfAdmin    = errors.New("attempting to give self admin")
	NotAnImage         = errors.New("data given was not an image")
)

func AddUsersApi(users *gin.RouterGroup) {
	users.GET("/all", func(c *gin.Context) {
		pageIndex, err := QueryParamUint(c, "pageIndex")
		if err != nil {
			var pI uint64 = 0
			pageIndex = &pI
		}
		users, err := database.SelectUsers(*pageIndex)
		JsonInfer(c, users, err)
	})

	users.POST("/update", func(c *gin.Context) {
		user, err := QueryParamToken(c)
		if JsonFail(c, err) {
			return
		}

		changed := false
		username, ok := c.GetQuery("username")
		if !ok {
			username = user.Username
			changed = true
		}
		email, ok := c.GetQuery("email")
		if !ok {
			email = user.Email
			changed = true
		}

		if !changed {
			JsonFail(c, EmptyUpdate)
			return
		}

		err = database.UpdateUser(user, username, email)
		JsonInfer(c, nil, err)
	})

	users.POST("/updateImage", func(c *gin.Context) {
		user, err := QueryParamToken(c)
		if JsonFail(c, err) {
			return
		}

		// @TODO Handle adversarial input
		if !strings.HasPrefix(c.ContentType(), "image/") {
			JsonFail(c, NotAnImage)
			return
		}

		// @TODO Handle adversarial input
		extension := strings.SplitN(c.ContentType(), "/", 2)[1]

		// @TODO Limit size of file
		image, err := ioutil.ReadAll(c.Request.Body)
		if JsonFail(c, err) {
			return
		}

		imageId, err := database.InsertImage(image, extension)
		if JsonFail(c, err) {
			return
		}

		err = database.UpdateProfilePic(user, imageId)
		JsonInfer(c, nil, err)
	})

	users.POST("/addAdmin", func(c *gin.Context) {
		user, err := QueryParamToken(c)
		if JsonFail(c, err) {
			return
		}

		id, err := QueryParamUint(c, "id")
		if JsonFail(c, err) {
			return
		}

		if *id == user.Id {
			JsonFail(c, GivingSelfAdmin)
			return
		}

		if user.UserGroup != models.AdminUser {
			JsonFail(c, MissingPermissions)
			return
		}

		err = database.UpdateUserGroup(*id, models.AdminUser)
		JsonInfer(c, nil, err)
	})

	users.POST("/removeAdmin", func(c *gin.Context) {
		user, err := QueryParamToken(c)
		if JsonFail(c, err) {
			return
		}

		id, err := QueryParamUint(c, "id")
		if JsonFail(c, err) {
			return
		}

		if user.UserGroup != models.AdminUser {
			JsonFail(c, MissingPermissions)
			return
		}

		err = database.UpdateUserGroup(*id, models.NormalUser)
		JsonInfer(c, nil, err)
	})

	users.POST("/add", func(c *gin.Context) {
		username, ok := c.GetQuery("username")
		if !ok {
			JsonFail(c, MissingUsername)
			return
		}
		email, ok := c.GetQuery("email")
		if !ok {
			JsonFail(c, MissingEmail)
			return
		}
		password, ok := c.GetQuery("password")
		if !ok {
			JsonFail(c, MissingPassword)
			return
		}

		err := database.InsertUser(username, email, password, models.NormalUser)
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
	permissions.POST("/add", func(c *gin.Context) {
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

		err = database.AddPermission(user, *target, permission)
		JsonInfer(c, nil, err)
	})

	permissions.POST("/remove", func(c *gin.Context) {
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
