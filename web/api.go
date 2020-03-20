package web

import (
	"encoding/json"
	"fmt"
	"github.com/A1Liu/webserver/database"
	"log"
	"net/http"
	"strconv"
)

func serialize(object interface{}) string {
	objectJson, err := json.Marshal(object)
	if err != nil { // If serialization of a type fails, we've made a mistake
		log.Fatal(err)
	}
	return string(objectJson)
}

func ListUsers(w http.ResponseWriter, _ *http.Request) {
	users, err := database.SelectUsers(database.GetDb(), 50, 0)
	if err != nil {
		fmt.Fprint(w, serialize(err))
	} else {
		fmt.Fprintln(w, serialize(users))
	}
}

func Clear(_ http.ResponseWriter, _ *http.Request) {
	database.Clear()
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	userGroup, err := strconv.ParseUint(r.URL.Query()["userGroup"][0], 10, 64)
	if err != nil {
		fmt.Fprint(w, serialize(err))
		return
	}
	database.InsertUser(database.GetDb(), r.URL.Query()["user"][0],
		r.URL.Query()["password"][0], userGroup)
}

func AddBook(w http.ResponseWriter, r *http.Request) {
}

