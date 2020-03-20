package web

import (
	"fmt"
	"github.com/A1Liu/webserver/database"
	"net/http"
	"strconv"
)

func ListUsers(w http.ResponseWriter, _ *http.Request) {
	users, err := database.SelectUsers(database.GetDb(), 50, 0)
	if err != nil {
		fmt.Fprintln(w, err)
	} else {
		fmt.Fprintln(w, users)
	}
}

func Clear(_ http.ResponseWriter, _ *http.Request) {
	database.Clear()
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	userGroup, err := strconv.ParseUint(r.URL.Query()["userGroup"][0], 10, 64)
	if err != nil {
		fmt.Fprint(w, err)
	}
	database.InsertUser(database.GetDb(), r.URL.Query()["user"][0],
		r.URL.Query()["password"][0], userGroup)
}
