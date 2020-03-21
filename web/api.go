package web

import (
	"errors"
	"github.com/A1Liu/webserver/database"
	"github.com/A1Liu/webserver/models"
	"github.com/gin-gonic/gin"
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
		JsonFail(c, err)
	} else {
		c.JSON(200, OkApiMessage{200, object})
	}
}

func JsonFail(c *gin.Context, err error) {
	c.JSON(400, ErrorApiMessage{400, err.Error()})
}

func GetQueryParamLogin(c *gin.Context) (*models.User, error) {
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

	return database.AuthorizeWithPassword(database.GetDb(), *login, *passwordNullable)
}

func GetQueryParamToken(c *gin.Context) (*models.User, error) {
	token, ok := c.GetQuery("token")
	if !ok {
		return nil, MissingToken
	} else {
		return database.AuthorizeWithToken(database.GetDb(), token)
	}
}
