package main

import (
	"encoding/json"
	"log"
	"strings"
)

const (
	ContentTypeHeaderName      = "Content-Type"
	ContentTypeHeaderValueJson = "application/json"
)

func doMocking(url, requestMethod string, body []byte, header map[string][]string) (string,string,map[string][]string ) {

	var responseBody, responseContentType string
	var responseHeaders map[string][]string

	log.Printf("entering doMocking with url %s method %s body %s",url,requestMethod, body)


	if header[ContentTypeHeaderName]!=nil && strings.EqualFold(ContentTypeHeaderValueJson,header[ContentTypeHeaderName][0]) {
		var requestBodyJson interface{}
		err := json.Unmarshal(body, &requestBodyJson)
		if err != nil {
			panic(err)
		}
		log.Println(requestBodyJson)

		matchingApiConfig := findMatchingApiConfig(url,requestMethod)
		log.Print(matchingApiConfig)
	} else{
		log.Print("invalid Content-Type header")
	}

	log.Print("exiting doMocking")
	return responseBody,responseContentType,responseHeaders
}

func findMatchingApiConfig(urlToSearch, requestMethod string) ApiConfig {
	var matchingApiConfig ApiConfig

	log.Printf("inside findMatchingApiConfig to find matching config for url %s requestMethod %s",urlToSearch,requestMethod)

	var apiConfigArr = getApiConfigArray()

	if apiConfigArr!=nil {

		for _,apiConfig := range apiConfigArr {

			if strings.EqualFold(apiConfig.Url,urlToSearch) && strings.EqualFold(apiConfig.Method,requestMethod) {
				log.Print("matching api config found with Id ",apiConfig.Id)
				return apiConfig
			}

		}
	}

	log.Print("exiting findMatchingApiConfig")

	//TODO check how to return nil from here
	return matchingApiConfig
}
