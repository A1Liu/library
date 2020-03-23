package utils

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type QueryMap map[string]string

type Response struct {
	Status        uint64
	StatusMessage string
	Body          string
}

func SendRequest(method, endpoint string, queryParams QueryMap, formValues url.Values) Response {
	var req http.Request
	reqURL, err := url.Parse("http://localhost:8080/api/v1" + endpoint)
	FailIf(err, "failed to parse endpoint")
	q := reqURL.Query()
	for k, v := range queryParams {
		q.Add(k, v)
	}
	reqURL.RawQuery = q.Encode()
	req.URL = reqURL
	req.Method = method
	req.PostForm = formValues
	resp, err := http.DefaultClient.Do(&req)
	FailIf(err, "Failed to get from API")
	defer resp.Body.Close()

	statusMessage := strings.SplitN(resp.Status, " ", 2)
	statusCode, err := strconv.ParseUint(statusMessage[0], 10, 64)
	FailIf(err, "failed to parse status code")

	data, err := ioutil.ReadAll(resp.Body)
	FailIf(err, "Failed to read response body")
	dataString := strings.TrimSpace(string(data))
	return Response{statusCode, statusMessage[1], dataString}
}
