package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

const (
	ContentTypeHeaderName      = "Content-Type"
	ContentTypeHeaderValueJson = "application/json"
)

func doMocking(url, requestMethod string, requestBody []byte,
		requestHeader map[string][]string) (string,map[string][]string) {

	var responseBody string
	var responseHeaders map[string][]string

	log.Printf("entering doMocking with url %s method %s body %s",url,requestMethod, requestBody)

	if requestHeader[ContentTypeHeaderName]!=nil &&
		strings.Contains(requestHeader[ContentTypeHeaderName][0],ContentTypeHeaderValueJson) {
		var requestBodyJson map[string]interface{}
		err := json.Unmarshal(requestBody, &requestBodyJson)
		if err != nil {
			panic(err)
		}
		log.Println("parsed request body json is ",requestBodyJson)

		matchingApiConfig := findMatchingApiConfig(url,requestMethod)
		if matchingApiConfig == nil {
			log.Printf("no matching config found for this api request")
			responseBody="no matching config found for this api request"
		} else {
			log.Printf("found matching api config with id %s value %v ",matchingApiConfig.Id,matchingApiConfig)
			return getMockedResponse(matchingApiConfig,requestBodyJson,requestHeader)
		}
	} else{
		log.Print("invalid Content-Type header",requestHeader[ContentTypeHeaderName])
		responseBody= fmt.Sprintf("%s %v","invalid Content-Type header",requestHeader[ContentTypeHeaderName])
	}
	return responseBody,responseHeaders
}

func getMockedResponse(apiConfig *ApiConfig, requestBodyJsonMap map[string]interface{},
		requestHeaderMap map[string][]string) (string, map[string][]string) {

	var responseBody string
	var responseHeaders map[string][]string

	responseBodyConfigJsonMap := apiConfig.ResponseConfig.ResponseJsonBody
	// set the values in response json map based on response config
	setResponseBodyMap(responseBodyConfigJsonMap, requestBodyJsonMap)
	responseBodyBytes,err := json.Marshal(responseBodyConfigJsonMap)
	if err==nil {
		responseBody = string(responseBodyBytes)
	}
	// set response headers
	responseHeaderConfigJsonMap := apiConfig.ResponseConfig.ResponseHeaders
	// set the values in response json map based on response config
	responseHeaders=setResponseHeaderMap(responseHeaderConfigJsonMap, requestHeaderMap)

	return responseBody, responseHeaders
}

