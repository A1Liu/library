package test

import (
	"github.com/A1Liu/library/utils"
	"log"
	"net/url"
)

func ShouldFail(method, endpoint string, queryParams utils.QueryMap, formValues url.Values) {
	resp := utils.SendRequest(method, endpoint, queryParams, formValues)
	if resp.Status != 400 {
		log.Fatal("ERROR: endpoint ", endpoint, " should have failed for invalid params ", queryParams)
	} else {
		log.Println("SUCCESS: endpoint", endpoint, "failed for invalid params", queryParams)
	}
}

func ShouldSucceedReturning(method, endpoint string, queryParams utils.QueryMap,
	formValues url.Values) *utils.Response {
	resp := utils.SendRequest(method, endpoint, queryParams, formValues)
	if resp.Status != 200 {
		log.Println(resp.Body)
		log.Fatal("ERROR: endpoint ", endpoint, " should not have failed for params ", queryParams)
		return nil
	} else {
		return &resp
	}
}

func ShouldSucceed(method, endpoint string, queryParams utils.QueryMap,
	formValues url.Values, returnValue string) {
	resp := ShouldSucceedReturning(method, endpoint, queryParams, formValues)

	if resp.Body != returnValue {
		log.Fatal("ERROR: return value of ", resp.Body,
			" doesn't match expected value for params ", queryParams)
	}
	log.Println("SUCCESS: endpoint", endpoint, "succeeded for params", queryParams)
}
