package test

import (
	"encoding/json"
	"github.com/A1Liu/webserver/database"
	"github.com/A1Liu/webserver/models"
	"github.com/A1Liu/webserver/utils"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var rootUserToken string = ""

func GetRootUserToken() string {
	if rootUserToken != "" {
		return rootUserToken
	}

	database.ConnectToDb()
	_, err := database.InsertUser("root", "root@gmail.com", "rootpass", models.AdminUser)
	utils.FailIf(err, "couldn't connect to database")

	resp := ShouldSucceedReturning(http.MethodGet, "/users/token", utils.QueryMap{
		"login":    "root",
		"password": "rootpass",
	}, url.Values{})

	rootUserToken = resp.Body[1 : len(resp.Body)-1]
	return rootUserToken
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
	json.Unmarshal([]byte(resp.Body), &users)

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
