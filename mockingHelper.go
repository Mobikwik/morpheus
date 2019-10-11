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

	// set response headers
	responseHeaderConfigJsonMap := apiConfig.ResponseConfig.ResponseHeaders
	// set the values in response json map based on response config
	responseHeaders=setResponseHeaderMap(responseHeaderConfigJsonMap, requestHeaderMap)

	if err==nil {
		responseBody = string(responseBodyBytes)
	}

	return responseBody, responseHeaders

}

func setResponseHeaderMap(responseHeaderConfigJsonMap map[string]interface{}, requestHeaderMap map[string][]string) map[string][]string {

	var responseHeaderValuesJsonMap = make(map[string][]string)
	for headerName, responseConfigKeyValueGenericType := range responseHeaderConfigJsonMap {

		log.Printf("setting value for response header %s of type %T config value %s",
			headerName,responseConfigKeyValueGenericType,responseConfigKeyValueGenericType)

		var responseHeaderValueArr []string

		switch responseConfigKeyValue:= responseConfigKeyValueGenericType.(type) {
		case string:
			responseHeaderValueArr= getResponseHeaderConfigValueFromRequestHeader(responseConfigKeyValue, requestHeaderMap)
			log.Printf("setting single value %s for header %s", responseHeaderValueArr,headerName)
		/*case []string:
			for i, responseConfigKeyValueSingle := range responseConfigKeyValue {
				responseHeaderValueArr = append(responseHeaderValueArr, getResponseHeaderConfigValueFromRequestHeader(responseConfigKeyValueSingle, requestHeaderMap)[0])
				log.Printf("adding array value %s on index %d for header %s ", responseHeaderValueArr[i],i,headerName)
			}*/
		case  []interface{}:
			log.Printf("handling response header value config for header %s config type is %T",headerName,responseConfigKeyValueGenericType)

			for i, responseConfigKeyValueSingle := range responseConfigKeyValue {
				log.Printf("getting value for config %s",responseConfigKeyValueSingle)
				responseConfigKeyValueSingleStr, ok := responseConfigKeyValueSingle.(string)
				if ok {
					responseHeaderValueArr = append(responseHeaderValueArr, getResponseHeaderConfigValueFromRequestHeader(responseConfigKeyValueSingleStr, requestHeaderMap)[0])
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

func getResponseHeaderConfigValueFromRequestHeader(responseHeaderConfigValue string, requestHeaderMap map[string][]string) []string {
	log.Printf("inside getResponseHeaderConfigValueFromRequestHeader, getting value for response config %s", responseHeaderConfigValue)
	if strings.HasPrefix(responseHeaderConfigValue, "requestHeaders.") {
		responseHeaderConfigValueSplit := strings.Split(responseHeaderConfigValue, ".")

		if len(responseHeaderConfigValueSplit) < 2 {
			// throw error
			log.Print("invalid response header configuration")
			return nil
		}
		/* golang converts request headers to canonical form, so we need to do the same while fetching header values
		https://godoc.org/net/http#CanonicalHeaderKey*/
		canonicalHeaderName := http.CanonicalHeaderKey(responseHeaderConfigValueSplit[1])
		/* if config is like $requestHeaders.Content-Type[2], get the array index part,i.e. 2
		responseHeaderConfigValueSplit[1] = Content-Type[2], len(responseHeaderConfigValueSplit[1]) =15,so:
		openingBracketIndex(index of [) = 12 = 15-3
		closingBracketIndex(index of ]) = 14 = 15-1
		arrIndex will have value 2
		*/
		openingBracketIndex := strings.Index(responseHeaderConfigValueSplit[1], "[")
		closingBracketIndex := strings.Index(responseHeaderConfigValueSplit[1], "]")

		if openingBracketIndex == len(responseHeaderConfigValueSplit[1])-3 && closingBracketIndex == len(responseHeaderConfigValueSplit[1])-1 {
			// remove [2] part from Content-Type[2], set canonicalHeaderName to Content-Type
			canonicalHeaderName=canonicalHeaderName[0:openingBracketIndex]
			arrIndex := responseHeaderConfigValueSplit[1][openingBracketIndex+1:closingBracketIndex]
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
		return []string{responseHeaderConfigValue}
	}
	return []string{}
}


func setResponseBodyMap(responseBodyConfigJsonMap map[string]interface{}, requestBodyJsonMap map[string]interface{}) {

	for key, responseBodyConfigValueGenericType := range responseBodyConfigJsonMap {
		log.Printf("getting value for key %s response body config value %v of type %T",key, responseBodyConfigValueGenericType,
			responseBodyConfigValueGenericType)

		var responseBodyValueArr []interface{}

		switch responseBodyConfigValue := responseBodyConfigValueGenericType.(type) {

		case []string:
			for i, responseConfigValueSingle := range responseBodyConfigValue {
				responseBodyValueArr = append(responseBodyValueArr, getResponseBodyValueFromRequestBody(responseConfigValueSingle,
					requestBodyJsonMap))
				log.Printf("adding array value %v on index %d for header %s ", responseBodyValueArr[i],i,key)
			}
			responseBodyConfigJsonMap[key]=responseBodyValueArr
		case []interface {}:
			for i, responseBodyConfigKeyValueSingle := range responseBodyConfigValue {
				log.Printf("getting value for config %s",responseBodyConfigKeyValueSingle)
				responseBodyConfigValueSingleStr, ok := responseBodyConfigKeyValueSingle.(string)
				if ok {
					responseBodyValueArr = append(responseBodyValueArr, getResponseBodyValueFromRequestBody(responseBodyConfigValueSingleStr, requestBodyJsonMap))
					log.Printf("adding array value %v on index %d for header %s ", responseBodyValueArr[i],i,key)
				}
			}
			responseBodyConfigJsonMap[key]=responseBodyValueArr
		case string:
			responseBodyConfigJsonMap[key] = getResponseBodyValueFromRequestBody(responseBodyConfigValue, requestBodyJsonMap)

			// when the value is a nested json, do recursive call
		case map[string]interface{}:
			setResponseBodyMap(responseBodyConfigValue, requestBodyJsonMap)
		default:
			fmt.Printf("no processing needed for response body config %v type %T",responseBodyConfigValueGenericType,
				responseBodyConfigValueGenericType)
		}
	}
}

func getResponseBodyValueFromRequestBody(responseBodyConfigValue string, requestBodyJsonMap map[string]interface{}) interface{} {
	log.Printf("inside getResponseBodyValueFromRequestBody, getting value for response config %s", responseBodyConfigValue)
	if strings.HasPrefix(responseBodyConfigValue, "requestJsonBody.") {

		responseBodyConfigValueSplit := strings.Split(responseBodyConfigValue, ".")

		// return nil for invalid config value "requestJsonBody."
		if len(responseBodyConfigValueSplit) < 2 {
			//TODO throw error
			log.Print("invalid response body configuration ", responseBodyConfigValue)
			return nil
		}

		// temp variables to hold values fetched from requestBodyJsonMap.
		// Declaring here because these can't be declared inside below for loop as they need to hold values of previous iteration
		var requestBodyValueMapOfInterfaceType map[string]interface{}
		var requestBodyValueInterfaceType interface{}
		ok1 := false

		/*	get the value of first nesting level object reference
			for ex: if config is requestJsonBody.orderDetails.addressDetails.pincode, get value of requestJsonBody["orderDetails"]
			and store in requestBodyValueInterfaceType
		   Value can be of any type (string,number or another nested object, so storing in interface{} type)
		*/
		requestBodyValueInterfaceType = requestBodyJsonMap[responseBodyConfigValueSplit[1]]

		/* process all the nested object references from 2nd level onwards by looping around the array split with seperator ".",
		   i.e. iteration i=1: get value of requestJsonBody.[orderDetails].[addressDetails] from requestBodyValueInterfaceType and store in requestBodyValueInterfaceType
				iteration i=2: get value of requestJsonBody.[orderDetails].[addressDetails].pincode from requestBodyValueInterfaceType and store in requestBodyValueInterfaceType
		*/
		for i := 1; i < len(responseBodyConfigValueSplit); i++ {

			/* requestBodyValueMapOfInterfaceType is typecast of requestBodyValueInterfaceType from interface{} type to map[string]interface{}
			 i.e. value of requestJsonBody.[orderDetails].[addressDetails] is a nested json, so storing this value in
			requestBodyValueMapOfInterfaceType to extract further values (like requestJsonBody.[orderDetails].[addressDetails].pincode) from it in next iteration
			*/
			// Don't use ":=" for value assignment as it will redeclare requestBodyValueMapOfInterfaceType as a new local variable in each iteration
			requestBodyValueMapOfInterfaceType, ok1 = requestBodyValueInterfaceType.(map[string]interface{})

			// checking if we have more nested config values at (i+1)th level
			if ok1 && i+1 < len(responseBodyConfigValueSplit) {

				jsonKeyName:= responseBodyConfigValueSplit[i+1]

				/* if config is like $requestJsonBody.txnTypes[2], get the array index part,i.e. 2
				responseBodyConfigValueSplit[1] = txnTypes[2], len(responseBodyConfigValueSplit[1]) =11,so:
				openingBracketIndex(index of [) = 8 = 11-3
				closingBracketIndex(index of ]) = 10 = 11-1
				arrIndex will have value 2
				*/
				openingBracketIndex := strings.Index(jsonKeyName, "[")
				closingBracketIndex := strings.Index(jsonKeyName, "]")

				if openingBracketIndex == len(jsonKeyName)-3 && closingBracketIndex == len(jsonKeyName)-1 {
					// get substring "2" from txnTypes[2]
					arrIndex := jsonKeyName[openingBracketIndex+1:closingBracketIndex]
					// get substring "txnTypes" from txnTypes[2] and set in jsonKeyName
					jsonKeyName=jsonKeyName[0:openingBracketIndex]

					arrIndexInt, err := strconv.Atoi(arrIndex) // convert string "2" to int
					if err!=nil {
						log.Printf("error parsing string to int for invalid config %s %v",responseBodyConfigValue,err)
						return nil // invalid config
					}
					var interfaceArrayTypeValue []interface{}
					ok2:=false
					/* As we are using config like txnTypes[2], value of txnTypes must be an array,
					so typecast to []interface{} */
					interfaceArrayTypeValue, ok2 = (requestBodyValueMapOfInterfaceType[jsonKeyName]).([]interface{})
					if ok2 {
						requestBodyValueInterfaceType = interfaceArrayTypeValue[arrIndexInt]
					} else {
						//throw error
						log.Printf("invalid config %s, trying to typecast non-array to array",responseBodyConfigValue)
						return nil
					}
				} else if openingBracketIndex == -1 && closingBracketIndex == -1 {
					// set value for next iteration if the config is like $requestHeaders.Content-Type i.e. does not have "[" and "]"
					requestBodyValueInterfaceType = requestBodyValueMapOfInterfaceType[jsonKeyName]
				}
			} else {
				/* reached end of nested values i.e. requestJsonBody.[orderDetails].[addressDetails].pincode, no more iterations possible,
				so return the final value from request body */
				return requestBodyValueInterfaceType
			}
		}
	} else {
		// response config does not start with "requestJsonBody.", so it's hard-coded value, so return that value as it is
		return responseBodyConfigValue
	}
	return nil
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
