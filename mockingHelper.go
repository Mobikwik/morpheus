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

func doMocking(url, requestMethod string, bodyString []byte, header map[string][]string) (string,string,map[string][]string ) {

	var responseBody, responseContentType string
	var responseHeaders map[string][]string

	log.Printf("entering doMocking with url %s method %s body %s",url,requestMethod,bodyString)


	if header[ContentTypeHeaderName]!=nil && strings.EqualFold(ContentTypeHeaderValueJson,header[ContentTypeHeaderName][0]) {
		var requestBodyJson interface{}
		err := json.Unmarshal(bodyString, &requestBodyJson)
		if err != nil {
			panic(err)
		}
		log.Println(requestBodyJson)
	} else{
		log.Print("invalid Content-Type header")
	}


	log.Print("exiting doMocking")
	return responseBody,responseContentType,responseHeaders
}
