package main

import (
	"github.com/A1Liu/webserver/test"
	"github.com/A1Liu/webserver/utils"
	"log"
	"net/http"
	"net/url"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	test.ShouldSucceed(http.MethodGet, "/clear", utils.QueryMap{}, url.Values{}, "null")

	token := test.GetRootUserToken()

	log.Println(token)
	test.TestUserAdd()
}
