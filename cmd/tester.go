package main

import (
	"github.com/A1Liu/library/test"
	"github.com/A1Liu/library/utils"
	"log"
	"net/http"
	"net/url"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	test.ShouldSucceed(http.MethodGet, "/clear", utils.QueryMap{}, url.Values{}, "null")

	test.TestUserPermissionsAdd()
	test.TestUserAdd()
}
