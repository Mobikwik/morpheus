package service

import (
	"encoding/json"
	"github.com/Mobikwik/morpheus/commons"
	"github.com/Mobikwik/morpheus/model"
	"log"
	"time"
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
	requestHeader map[string][]string) (string, map[string][]string, int) {

	/*	// this function will be called in case of any "panic"
		defer recoverConfigError2()*/

	var responseBody string
	var responseHeaders map[string][]string

	log.Printf("entering doMocking with url %s method %s request header %v request body %s",
		url, requestMethod, requestHeader, requestBody)

	//TODO put a check on request method. A POST api should give error if called via GET

	// remove the content-type header check
	/*if requestHeader[ContentTypeHeaderName] != nil &&
	strings.Contains(requestHeader[ContentTypeHeaderName][0], ContentTypeHeaderValueJson) {*/
	var requestBodyMap map[string]interface{}
	err := json.Unmarshal(requestBody, &requestBodyMap)
	if err != nil {
		panic(err)
	}
	log.Println("parsed request body json is ", requestBodyMap)

	matchingMockConfig := commons.FindMatchingMockConfig(url, requestMethod, requestBodyMap)
	if matchingMockConfig == nil {
		log.Printf("no matching config found for this api request")
		responseBody = "no matching config found for this api request"
	} else {
		responseBody, responseHeaders := getMockedResponse(matchingMockConfig, requestBodyMap, requestHeader)

		// check if api config has any setting for introducing delay in sending response. This is to test api timeouts
		if matchingMockConfig.ResponseDelayInSeconds > 0 {
			// time.Duration by default is in nanoseconds, converting it in seconds
			var responseDelay = time.Duration(matchingMockConfig.ResponseDelayInSeconds) * time.Second
			log.Printf("introducing response delay of %s seconds", responseDelay)
			time.Sleep(responseDelay)
		}
		return responseBody, responseHeaders, matchingMockConfig.ResponseMockValues.HttpCode
	}
	/*} else {
		log.Print("invalid Content-Type header", requestHeader[ContentTypeHeaderName])
		responseBody = fmt.Sprintf("%s %v", "invalid Content-Type header", requestHeader[ContentTypeHeaderName])
	}*/
	return responseBody, responseHeaders, 200
}

func getMockedResponse(mockConfig *model.MockConfig, requestBodyJsonMap map[string]interface{},
	requestHeaderMap map[string][]string) (responseBody string, responseHeaders map[string][]string) {

	responseBodyConfigJsonMap := mockConfig.ResponseMockValues.ResponseBodyMockValues
	// set the values in response json map based on response config
	setResponseBodyMap(responseBodyConfigJsonMap, requestBodyJsonMap)
	responseBodyBytes, err := json.Marshal(responseBodyConfigJsonMap)
	if err == nil {
		responseBody = string(responseBodyBytes)
	}
	// set response headers
	responseHeaderConfigJsonMap := mockConfig.ResponseMockValues.ResponseHeadersMockValues
	// set the values in response json map based on response config
	responseHeaders = setResponseHeaderMap(responseHeaderConfigJsonMap, requestHeaderMap)

	return responseBody, responseHeaders
}
