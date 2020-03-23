package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type QueryMap map[string]string

func failIf(err error, msg string) {
	if err != nil {
		log.Fatal("ERROR: ", msg)
	}
}

func shouldFail(endpoint string, queryParams QueryMap) {
	var req http.Request
	reqURL, err := url.Parse("http://localhost:8080/api/v1" + endpoint)
	failIf(err, "failed to parse endpoint")
	q := reqURL.Query()
	for k, v := range queryParams {
		q.Add(k, v)
	}
	reqURL.RawQuery = q.Encode()
	req.URL = reqURL
	req.Method = http.MethodGet
	resp, err := http.DefaultClient.Do(&req)
	failIf(err, "Failed to get from API")

	defer resp.Body.Close()

	if resp.Status != "400 Bad Request" {
		log.Fatal("ERROR: Endpoint ", endpoint, " should have failed for invalid params ", queryParams)
	} else {
		log.Println("SUCESS: Endpoint", endpoint, "failed for invalid params", queryParams)
	}
}

func shouldSucceedReturning(endpoint string, queryParams QueryMap) string {
	var req http.Request
	reqURL, err := url.Parse("http://localhost:8080/api/v1" + endpoint)
	failIf(err, "failed to parse endpoint")
	q := reqURL.Query()
	for k, v := range queryParams {
		q.Add(k, v)
	}
	reqURL.RawQuery = q.Encode()
	req.URL = reqURL
	req.Method = http.MethodGet
	resp, err := http.DefaultClient.Do(&req)
	failIf(err, "Failed to get from API")

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	dataString := strings.TrimSpace(string(data))
	if resp.Status != "200 OK" {
		log.Println(dataString)
		log.Fatal("ERROR: Endpoint ", endpoint, " should not have failed for params ", queryParams)
		return ""
	} else {
		failIf(err, "Failed to read response from API")
		return dataString
	}
}

func shouldSucceed(endpoint string, queryParams QueryMap, returnValue string) {
	dataString := shouldSucceedReturning(endpoint, queryParams)

	if dataString != returnValue {
		log.Fatal("ERROR: Return value of ", dataString, " doesn't match expected value for params ", queryParams)
	}
	log.Println("SUCESS: Endpoint", endpoint, "succeeded for params", queryParams)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	shouldSucceed("/clear", QueryMap{}, "null")

	shouldFail("/users/add", QueryMap{})
	shouldSucceed("/users/all", QueryMap{}, "[]")

	_ = shouldSucceedReturning("/users/add", QueryMap{
		"username": "hi",
		"email":    "hello@gmail.com",
		"password": "asdfghjkl",
	})

}
