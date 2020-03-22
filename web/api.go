package web

import (
	"errors"
	"fmt"
	"github.com/A1Liu/webserver/database"
	"github.com/A1Liu/webserver/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ErrorApiMessage struct {
	Status  uint64 `json:"status"`
	Message string `json:"message"`
}

type OkApiMessage struct {
	Status uint64      `json:"status"`
	Value  interface{} `json:"value"`
}

var (
	MissingLogin       = errors.New("missing login query parameter")
	MissingPassword    = errors.New("missing password query parameter")
	NoLoginInformation = errors.New("neither login nor password was provided")
	MissingToken       = errors.New("missing token")
	TooManyAuthMethods = errors.New("gave too many authorization methods")
)

func JsonInfer(c *gin.Context, object interface{}, err error) {
	if err != nil {
		c.JSON(400, ErrorApiMessage{400, err.Error()})
	} else {
		c.JSON(200, OkApiMessage{200, object})
	}
}

func JsonFail(c *gin.Context, err error) bool {
	if err != nil {
		c.JSON(400, ErrorApiMessage{400, err.Error()})
		return true
	}
	return false
}

func QueryParamLogin(c *gin.Context) (*models.User, error) {
	var login *string
	usernameOrEmail, ok := c.GetQuery("login")
	if ok {
		login = &usernameOrEmail
	} else {
		login = nil
	}

	var passwordNullable *string
	password, ok := c.GetQuery("password")
	if ok {
		passwordNullable = &password
	} else {
		passwordNullable = nil
	}

	if login == nil && passwordNullable != nil {
		return nil, MissingLogin
	} else if login != nil && passwordNullable == nil {
		return nil, MissingPassword
	} else if login == nil && passwordNullable == nil {
		return nil, NoLoginInformation
	}

	return database.AuthorizeWithPassword(*login, *passwordNullable)
}

func QueryParamToken(c *gin.Context) (*models.User, error) {
	token, ok := c.GetQuery("token")
	if !ok {
		return nil, MissingToken
	} else {
		return database.AuthorizeWithToken(token)
	}
}

func QueryParamUint(c *gin.Context, param string) (*uint64, error) {
	valString, ok := c.GetQuery(param)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Missing required param `%s`", param))
	}

	val, err := strconv.ParseUint(valString, 10, 64)
	return &val, err
}
