package test

import (
	"encoding/json"
	"github.com/A1Liu/library/database"
	"github.com/A1Liu/library/models"
	"github.com/A1Liu/library/utils"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	rootUsername = "root"
	rootEmail    = "root@gmail.com"
	rootPassword = "rootpass"
)

var rootUserToken string = ""

func GetRootUserToken() string {
	if rootUserToken != "" {
		return rootUserToken
	}

	database.ConnectToDb()
	err := database.InsertUser(rootUsername, rootEmail, rootPassword, models.AdminUser)
	utils.FailIf(err, "couldn't connect to database")

	resp := ShouldSucceedReturning(http.MethodGet, "/users/token", utils.QueryMap{
		"login":    rootUsername,
		"password": rootPassword,
	}, url.Values{})

	rootUserToken = resp.Body[1 : len(resp.Body)-1]
	return rootUserToken
}

func TestUserPermissionsAdd() {
	resp := ShouldSucceedReturning(http.MethodGet, "/users/get", utils.QueryMap{
		"token": GetRootUserToken(),
	}, url.Values{})

	var user models.User
	err := json.Unmarshal([]byte(resp.Body), &user)
	utils.FailIf(err, "json unmarshalling failed for /users/get")

	if user.Username != rootUsername {
		utils.Fail("Username incorrect")
	}
	if user.Email != rootEmail {
		utils.Fail("Email incorrect")
	}
	if user.UserGroup != models.AdminUser {
		utils.Fail("usergroup incorrect")
	}
}

func TestUserAdd() {
	ShouldFail(http.MethodGet, "/users/add", utils.QueryMap{}, url.Values{})

	username, email, password := "hi", "hello@gmail.com", "asdfghjkl"
	_ = ShouldSucceedReturning(http.MethodGet, "/users/add", utils.QueryMap{
		"username": username,
		"email":    email,
		"password": password,
	}, url.Values{})

	var users []models.User
	resp := ShouldSucceedReturning(http.MethodGet, "/users/all", utils.QueryMap{}, url.Values{})
	err := json.Unmarshal([]byte(resp.Body), &users)
	utils.FailIf(err, "json unmarshalling failed for /users/all")

	for _, user := range users {
		if user.Username != strings.ToLower(username) || user.Email != email ||
			user.UserGroup != models.NormalUser {
			continue
		}

		log.Println("SUCCESS: endpoint /users/all contains the user we created")

		return
	}
	utils.Fail("couldn't find user")
}
