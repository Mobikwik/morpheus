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

func doMocking(url, requestMethod string, requestBody []byte, requestHeader map[string][]string) (string,string,map[string][]string ) {

	var responseBody, responseContentType string
	var responseHeaders map[string][]string

	log.Printf("entering doMocking with url %s method %s body %s",url,requestMethod, requestBody)


	if requestHeader[ContentTypeHeaderName]!=nil &&
		strings.EqualFold(ContentTypeHeaderValueJson,requestHeader[ContentTypeHeaderName][0]) {
		var requestBodyJson map[string]interface{}
		err := json.Unmarshal(requestBody, &requestBodyJson)
		if err != nil {
			panic(err)
		}
		log.Println(requestBodyJson)

		matchingApiConfig := findMatchingApiConfig(url,requestMethod)
		log.Print(matchingApiConfig)

		if matchingApiConfig == nil {
			// return nil values or throw error
		} else {
			return getMockedResponse(matchingApiConfig,requestBodyJson,requestHeader)
		}
	} else{
		log.Print("invalid Content-Type header")
	}

	log.Print("exiting doMocking")
	return responseBody,responseContentType,responseHeaders
}

func getMockedResponse(apiConfig *ApiConfig, requestBodyJsonMap map[string]interface{}, requestHeaderMap map[string][]string) (string, string, map[string][]string) {

	var responseBody, responseContentType string
	var responseHeaders map[string][]string

	responseBodyConfigJsonMap := apiConfig.ResponseConfig.ResponseJsonBody
	// set the values in response json map based on response config
	setResponseBodyMap(responseBodyConfigJsonMap, requestBodyJsonMap)

	// set response headers
	responseHeaderConfigJsonMap := apiConfig.ResponseConfig.ResponseHeaders
	// set the values in response json map based on response config
	setResponseHeaderMap(responseHeaderConfigJsonMap, requestHeaderMap)

	responseBodyBytes,err := json.Marshal(responseBodyConfigJsonMap)
	if err==nil {
		responseBody = string(responseBodyBytes)
	}

	return responseBody, responseContentType, responseHeaders

}

func setResponseHeaderMap(responseHeaderConfigJsonMap map[string]string, requestHeaderMap map[string][]string) {

	for key, responseKeyValue := range responseHeaderConfigJsonMap {
		fmt.Println("key:", key, "v:", responseKeyValue)

		if strings.HasPrefix(responseKeyValue, "requestHeaderBody.") {
			strValueSplit := strings.Split(responseKeyValue, ".")

			if len(strValueSplit) < 2 {
				// throw error
				log.Print("invalid response header configuration")
			}
			responseHeaderConfigJsonMap[key]=requestHeaderMap[strValueSplit[1]][0]
		}
	}
}


func setResponseBodyMap(responseBodyConfigJsonMap map[string]interface{}, requestBodyJsonMap map[string]interface{}) {
	for key, responseKeyValueGenericType := range responseBodyConfigJsonMap {
		fmt.Println("key:", key, "v:", responseKeyValueGenericType)

		switch responseKeyValue := responseKeyValueGenericType.(type) {

		case string:

			if strings.HasPrefix(responseKeyValue, "requestJsonBody.") {
				strValueSplit := strings.Split(responseKeyValue, ".")

				if len(strValueSplit) < 2 {
					// throw error
					log.Print("invalid response body configuration")
				}

				var mapTypeValue map[string]interface{} // temp variable to hold map type value
				ok := false
				/*
					get the value of first nesting level
					for ex: if config is requestJsonBody.orderDetails.orderId, get value of requestJsonBody["orderDetails"]
				*/
				interfaceTypeValue := requestBodyJsonMap[strValueSplit[1]]

				for i := 1; i < len(strValueSplit); i++ {

					mapTypeValue, ok = interfaceTypeValue.(map[string]interface{})

					if ok && i+1 < len(strValueSplit) {
						interfaceTypeValue = mapTypeValue[strValueSplit[i+1]]
					} else {
						responseBodyConfigJsonMap[key] = interfaceTypeValue
					}
				}
			}

		case map[string]interface{}:
			// recursive call to process the map
			setResponseBodyMap(responseKeyValue, requestBodyJsonMap)
		default:
			fmt.Printf("type is %T responseKeyValueGenericType %s \n", responseKeyValue, responseKeyValueGenericType)
		}
	}
}

func processStringValueOfResponseConfig(responseKeyValue string, requestBodyJsonMap map[string]interface{},
			responseBodyConfigJsonMap map[string]interface{}, key string) {
	if strings.HasPrefix(responseKeyValue, "requestJsonBody.") {
		strValueSplit := strings.Split(responseKeyValue, ".")

		if len(strValueSplit) < 2 {
			// throw error
			log.Print("invalid response responseKeyValueGenericType configuration")
		}

		var mapTypeValue map[string]interface{} // temp variable to hold map type value
		ok := false
		/*
			get the value of first nesting level
			for ex: if config is requestJsonBody.orderDetails.orderId, get value of requestJsonBody["orderDetails"]
		*/
		interfaceTypeValue := requestBodyJsonMap[strValueSplit[1]]

		for i := 1; i < len(strValueSplit); i++ {

			mapTypeValue, ok = interfaceTypeValue.(map[string]interface{})

			if ok && i+1 < len(strValueSplit) {
				interfaceTypeValue = mapTypeValue[strValueSplit[i+1]]
			} else {
				responseBodyConfigJsonMap[key] = interfaceTypeValue
			}
		}
	}
}

func findMatchingApiConfig(urlToSearch, requestMethod string) *ApiConfig {
	//var matchingApiConfig *ApiConfig

	log.Printf("inside findMatchingApiConfig to find matching config for url %s requestMethod %s",urlToSearch,requestMethod)

	var apiConfigArr = getApiConfigArray()

	if apiConfigArr!=nil {

		for _,apiConfig := range apiConfigArr {

			if strings.EqualFold(apiConfig.Url,urlToSearch) && strings.EqualFold(apiConfig.Method,requestMethod) {
				log.Print("matching api config found with Id ",apiConfig.Id)
				return &apiConfig
			}

		}
	}

	log.Print("exiting findMatchingApiConfig")

	return nil
}
