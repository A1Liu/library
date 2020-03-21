package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/A1Liu/webserver/database"
	"github.com/A1Liu/webserver/models"
	"log"
	"net/http"
	"strconv"
)

type ErrorApiMessage struct {
	Status uint64 `json:"status"`
	Message string `json:"message"`
}

type OkApiMessage struct {
	Status uint64 `json:"status"`
	Value interface {} `json:"value"`
}

var (
	MissingLogin = errors.New("missing login query parameter")
	MissingPassword = errors.New("missing password query parameter")
	MissingAuthorization = errors.New("need an authorization parameter (`token` or both of `login` and `password`)")
)

func serialize(w http.ResponseWriter, object interface{}) {
	var objectJson []byte
	var err error
	if e, ok := object.(error); ok {
		objectJson, err = json.Marshal(ErrorApiMessage{400, e.Error()})
	} else {
		objectJson, err = json.Marshal(OkApiMessage{200, object})
	}
	_, _ = fmt.Fprint(w, string(objectJson))
	if err != nil { // If serialization of a type fails, we've made a mistake
		log.Fatal(err)
	}
}

func getQueryParam(r *http.Request, param string) (string, error) {
	paramList :=  r.URL.Query()[param]
	if len(paramList) == 0 {
		return "", errors.New("missing param `" + param + "`")
	}
	return paramList[0], nil
}

func getQueryParamUint(r *http.Request, param string) (uint64, error) {
	paramList :=  r.URL.Query()[param]
	if len(paramList) == 0 {
		return 0, errors.New("missing param `" + param + "`")
	}
	group, err := strconv.ParseUint(paramList[0], 10, 64)
	return group, err
}

func getQueryParamLogin(r *http.Request) (*models.User, error) {
	var login *string
	usernameOrEmail, err := getQueryParam(r, "login")
	if err != nil {
		login = nil
	} else {
		login = &usernameOrEmail
	}

	var passwordNullable *string
	password, err := getQueryParam(r, "password")
	if err != nil {
		passwordNullable = nil
	} else {
		passwordNullable = &password
	}

	if login == nil && passwordNullable != nil {
		return nil, MissingLogin
	} else if login != nil && passwordNullable == nil {
		return nil, MissingPassword
	}

	token, err := getQueryParam(r, "password")
	if err != nil && login == nil {
		return nil, MissingAuthorization
	} else if login != nil {
		return database.AuthorizeWithPassword(database.GetDb(), *login, *passwordNullable)
	} else {
		return database.AuthorizeWithToken(token)
	}
}

func ListUsers(w http.ResponseWriter, _ *http.Request) {
	users, err := database.SelectUsers(database.GetDb(), 50, 0)
	if err != nil {
		serialize(w, err)
	} else {
		serialize(w, users)
	}
}

func Clear(w http.ResponseWriter, _ *http.Request) {
	serialize(w, database.Clear())
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	userGroup, err := getQueryParamUint(r, "userGroup")
	if err != nil {
		serialize(w, err)
		return
	}

	username, err := getQueryParam(r, "username")
	if err != nil {
		serialize(w, err)
		return
	}

	email, err := getQueryParam(r, "email")
	if err != nil {
		serialize(w, err)
		return
	}

	password, err := getQueryParam(r, "password")
	if err != nil {
		serialize(w, err)
		return
	}

	err = database.InsertUser(database.GetDb(), username, email, password, userGroup)
	serialize(w, err)
}

func AddBook(w http.ResponseWriter, r *http.Request) {
	title, err := getQueryParam(r, "title")
	if err != nil {
		serialize(w, err)
		return
	}

	description, err := getQueryParam(r, "description")
	if err != nil {
		serialize(w, err)
		return
	}

	user, err := getQueryParamLogin(r)
	if err != nil && err != MissingAuthorization {
		serialize(w, err)
		return
	}
	var id *uint64
	if user == nil {
		id = nil
	} else {
		id = &user.Id
	}

	serialize(w, database.InsertBook(database.GetDb(), id, title, description))
}

