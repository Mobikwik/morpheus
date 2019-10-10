package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	ContentTypeHeaderName      = "Content-Type"
	ContentTypeHeaderValueJson = "application/json"
)

func doMocking(url, requestMethod string, requestBody []byte,
		requestHeader map[string][]string) (string,string,map[string][]string) {

	var responseBody, responseContentType string
	var responseHeaders map[string][]string

	log.Printf("entering doMocking with url %s method %s body %s",url,requestMethod, requestBody)


	if requestHeader[ContentTypeHeaderName]!=nil &&
		strings.Contains(requestHeader[ContentTypeHeaderName][0],ContentTypeHeaderValueJson) {
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
	return responseBody,responseContentType,responseHeaders
}

func getMockedResponse(apiConfig *ApiConfig, requestBodyJsonMap map[string]interface{},
		requestHeaderMap map[string][]string) (string, string, map[string][]string) {

	var responseBody, responseContentType string
	var responseHeaders map[string][]string

	responseBodyConfigJsonMap := apiConfig.ResponseConfig.ResponseJsonBody
	// set the values in response json map based on response config
	setResponseBodyMap(responseBodyConfigJsonMap, requestBodyJsonMap)
	responseBodyBytes,err := json.Marshal(responseBodyConfigJsonMap)

	// set response headers
	responseHeaderConfigJsonMap := apiConfig.ResponseConfig.ResponseHeaders
	// set the values in response json map based on response config
	responseHeaders=setResponseHeaderMap(responseHeaderConfigJsonMap, requestHeaderMap)

	if err==nil {
		responseBody = string(responseBodyBytes)
	}

	return responseBody, responseContentType, responseHeaders

}

func setResponseHeaderMap(responseHeaderConfigJsonMap map[string]interface{}, requestHeaderMap map[string][]string) map[string][]string {

	var responseHeaderValuesJsonMap = make(map[string][]string)
	for headerName, responseConfigKeyValueGenericType := range responseHeaderConfigJsonMap {

		log.Printf("setting value for response header %s of type %T config value %s",
			headerName,responseConfigKeyValueGenericType,responseConfigKeyValueGenericType)

		var responseHeaderValueArr []string

		switch responseConfigKeyValue:= responseConfigKeyValueGenericType.(type) {
		case string:
			responseHeaderValueArr=getValueFromRequestHeader(responseConfigKeyValue, requestHeaderMap)
			log.Printf("setting single value %s for header %s", responseHeaderValueArr,headerName)
		case []string:
			for i, responseConfigKeyValueSingle := range responseConfigKeyValue {
				responseHeaderValueArr = append(responseHeaderValueArr, getValueFromRequestHeader(responseConfigKeyValueSingle, requestHeaderMap)[0])
				log.Printf("adding array value %s on index %d for header %s ", responseHeaderValueArr[i],i,headerName)
			}
		case  []interface{}:
			log.Printf("handling response header value config for header %s config type is %T",headerName,responseConfigKeyValueGenericType)

			for i, responseConfigKeyValueSingle := range responseConfigKeyValue {
				log.Printf("getting value for config %s",responseConfigKeyValueSingle)
				responseConfigKeyValueSingleStr, ok := responseConfigKeyValueSingle.(string)
				if ok {
					responseHeaderValueArr = append(responseHeaderValueArr, getValueFromRequestHeader(responseConfigKeyValueSingleStr, requestHeaderMap)[0])
					log.Printf("adding array value %s on index %d for header %s ", responseHeaderValueArr[i],i,headerName)
				}

			}

		default:
			log.Printf("invalid response header value config for header %s config type is %T",headerName,responseConfigKeyValueGenericType)

		}

		log.Printf("setting final value for response header %s = %s",headerName, responseHeaderValueArr)
		//responseHeaderConfigJsonMap[headerName] = responseHeaderValueArr
		responseHeaderValuesJsonMap[headerName] = responseHeaderValueArr
	}
	return responseHeaderValuesJsonMap
}

func getValueFromRequestHeader(responseConfigKeyValue string, requestHeaderMap map[string][]string) []string {
	log.Printf("inside getValueFromRequestHeader, getting value for response config %s",responseConfigKeyValue)
	if strings.HasPrefix(responseConfigKeyValue, "requestHeaders.") {
		strValueSplit := strings.Split(responseConfigKeyValue, ".")

		if len(strValueSplit) < 2 {
			// throw error
			log.Print("invalid response header configuration")
		}
		/* golang converts request headers to canonical form, so we need to do the same while fetching header values
		https://godoc.org/net/http#CanonicalHeaderKey*/
		canonicalHeaderName := http.CanonicalHeaderKey(strValueSplit[1])
		/* if config is like $requestHeaders.Content-Type[2], get the array index part,i.e. 2
		strValueSplit[1] = Content-Type[2], len(strValueSplit[1]) =15,so:
		openingBracketIndex(index of [) = 12 = 15-3
		closingBracketIndex(index of ]) = 14 = 15-1
		arrIndex will have value 2
		*/
		openingBracketIndex := strings.Index(strValueSplit[1], "[")
		closingBracketIndex := strings.Index(strValueSplit[1], "]")

		if openingBracketIndex == len(strValueSplit[1])-3 && closingBracketIndex == len(strValueSplit[1])-1 {
			// remove [2] part from Content-Type[2], set canonicalHeaderName to Content-Type
			canonicalHeaderName=canonicalHeaderName[0:openingBracketIndex]
			arrIndex := strValueSplit[1][openingBracketIndex+1:closingBracketIndex]
			arrIndexInt, err := strconv.Atoi(arrIndex) // convert string to int
			if err!=nil {
				log.Print("error parsing string to int ",err)
				return []string{}
			}
			return []string{requestHeaderMap[canonicalHeaderName][arrIndexInt]}
		} else {
			// if the config is like $requestHeaders.Content-Type,return entire array of this header value
			return requestHeaderMap[canonicalHeaderName]
		}
	} else {
		// return the value inside config as it is
		return []string{responseConfigKeyValue}
	}
	return []string{}
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
