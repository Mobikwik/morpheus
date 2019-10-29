package service

import (
	"encoding/json"
	"fmt"
	"github.com/Mobikwik/morpheus/commons"
	"github.com/Mobikwik/morpheus/model"
	"log"
	"strings"
	"time"
)

const (
	ContentTypeHeaderName      = "Content-Type"
	ContentTypeHeaderValueJson = "application/json"
)

/*
func recoverConfigError2() string {
	if r := recover(); r!= nil {
		log.Print("recovered from ", r)
	}
	return "recover return"
}
*/
func DoMocking(url, requestMethod string, requestBody []byte,
	requestHeader map[string][]string) (string, map[string][]string) {

	/*	// this function will be called in case of any "panic"
		defer recoverConfigError2()*/

	var responseBody string
	var responseHeaders map[string][]string

	log.Printf("entering doMocking with url %s method %s body %s", url, requestMethod, requestBody)

	if requestHeader[ContentTypeHeaderName] != nil &&
		strings.Contains(requestHeader[ContentTypeHeaderName][0], ContentTypeHeaderValueJson) {
		var requestBodyJson map[string]interface{}
		err := json.Unmarshal(requestBody, &requestBodyJson)
		if err != nil {
			panic(err)
		}
		log.Println("parsed request body json is ", requestBodyJson)

		matchingApiConfig := commons.FindMatchingApiConfig(url, requestMethod)
		if matchingApiConfig == nil {
			log.Printf("no matching config found for this api request")
			responseBody = "no matching config found for this api request"
		} else {
			log.Printf("found matching api config with id %v value %v", matchingApiConfig.Id, matchingApiConfig)
			responseBody, responseHeaders := getMockedResponse(matchingApiConfig, requestBodyJson, requestHeader)

			// check if api config has any setting for introducing delay in sending response. This is to test api timeouts
			if matchingApiConfig.ResponseDelayInSeconds > 0 {
				// time.Duration by default is in nanoseconds, converting it in seconds
				var responseDelay = time.Duration(matchingApiConfig.ResponseDelayInSeconds) * time.Second
				log.Printf("introducing response delay of %s seconds", responseDelay)
				time.Sleep(responseDelay)
			}
			return responseBody, responseHeaders
		}
	} else {
		log.Print("invalid Content-Type header", requestHeader[ContentTypeHeaderName])
		responseBody = fmt.Sprintf("%s %v", "invalid Content-Type header", requestHeader[ContentTypeHeaderName])
	}
	return responseBody, responseHeaders
}

func getMockedResponse(apiConfig *model.ApiConfig, requestBodyJsonMap map[string]interface{},
	requestHeaderMap map[string][]string) (string, map[string][]string) {

	var responseBody string
	var responseHeaders map[string][]string

	responseBodyConfigJsonMap := apiConfig.ResponseConfig.ResponseJsonBody
	// set the values in response json map based on response config
	setResponseBodyMap(responseBodyConfigJsonMap, requestBodyJsonMap)
	responseBodyBytes, err := json.Marshal(responseBodyConfigJsonMap)
	if err == nil {
		responseBody = string(responseBodyBytes)
	}
	// set response headers
	responseHeaderConfigJsonMap := apiConfig.ResponseConfig.ResponseHeaders
	// set the values in response json map based on response config
	responseHeaders = setResponseHeaderMap(responseHeaderConfigJsonMap, requestHeaderMap)

	return responseBody, responseHeaders
}
