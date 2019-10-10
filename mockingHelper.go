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

func getResponseHeaderConfigValueFromRequestHeader(responseConfigKeyValue string, requestHeaderMap map[string][]string) []string {
	log.Printf("inside getResponseHeaderConfigValueFromRequestHeader, getting value for response config %s",responseConfigKeyValue)
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
	for key, responseBodyKeyValueGenericType := range responseBodyConfigJsonMap {
		log.Printf("getting value for key %s response body config value %s of type %T",key,responseBodyKeyValueGenericType,
			responseBodyKeyValueGenericType)

		var responseBodyValueArr []interface{}

		switch responseBodyConfigKeyValue := responseBodyKeyValueGenericType.(type) {

		case []string:
			for i, responseConfigKeyValueSingle := range responseBodyConfigKeyValue {
				responseBodyValueArr = append(responseBodyValueArr,getResponseBodyConfigValueFromRequestBody(responseConfigKeyValueSingle, requestBodyJsonMap))
				log.Printf("adding array value %s on index %d for header %s ", responseBodyValueArr[i],i,key)
			}
			responseBodyConfigJsonMap[key]=responseBodyValueArr
		case []interface {}:

			for i, responseBodyConfigKeyValueSingle := range responseBodyConfigKeyValue {

				log.Printf("getting value for config %s",responseBodyConfigKeyValueSingle)
				responseBodyConfigKeyValueSingleStr, ok := responseBodyConfigKeyValueSingle.(string)
				if ok {
					responseBodyValueArr = append(responseBodyValueArr,getResponseBodyConfigValueFromRequestBody(responseBodyConfigKeyValueSingleStr, requestBodyJsonMap))
					log.Printf("adding array value %v on index %d for header %s ", responseBodyValueArr[i],i,key)
				}

			}
			responseBodyConfigJsonMap[key]=responseBodyValueArr
		case string:
			responseBodyConfigJsonMap[key] = getResponseBodyConfigValueFromRequestBody(responseBodyConfigKeyValue, requestBodyJsonMap)

			// when the value is a nested json, do recursive call
		case map[string]interface{}:
			setResponseBodyMap(responseBodyConfigKeyValue, requestBodyJsonMap)
		default:
			fmt.Printf("type is %T responseKeyValueGenericType %s \n", responseBodyConfigKeyValue, responseBodyKeyValueGenericType)
		}
	}
}

func getResponseBodyConfigValueFromRequestBody(responseBodyConfigKeyValue string, requestBodyJsonMap map[string]interface{}) interface{} {
	log.Printf("inside getResponseBodyConfigValueFromRequestBody, getting value for response config %s",responseBodyConfigKeyValue)
	if strings.HasPrefix(responseBodyConfigKeyValue, "requestJsonBody.") {

		strValueSplit := strings.Split(responseBodyConfigKeyValue, ".")

		if len(strValueSplit) < 2 {
			// throw error
			log.Print("invalid response body configuration")
		}


		var mapInterfaceTypeValue map[string]interface{} // temp variable to hold map type value
		ok1 := false
		ok2:=false
		/*
			get the value of first nesting level object reference
			for ex: if config is requestJsonBody.orderDetails.addressDetails.pincode, get value of requestJsonBody["orderDetails"]
			and store in interfaceTypeValue
		*/
		interfaceTypeValue := requestBodyJsonMap[strValueSplit[1]]
		/* process all the nested object references from 2nd level onwards by looping around the array split with ".",
		   i.e. iteration i=1: get value of requestJsonBody.[orderDetails].[addressDetails] from interfaceTypeValue and store in interfaceTypeValue
				iteration i=2: get value of requestJsonBody.[orderDetails].[addressDetails].pincode from interfaceTypeValue and store in interfaceTypeValue
		*/
		for i := 1; i < len(strValueSplit); i++ {

			/* mapTypeValue is typecast of interfaceTypeValue from interface{} type to map[string]interface{}
			 i.e. value of requestJsonBody.[orderDetails].[addressDetails] is a nested json, so storing this value in
			mapTypeValue to extract further values (like requestJsonBody.[orderDetails].[addressDetails].pincode) from it in next iteration
			*/
			mapInterfaceTypeValue, ok1 = interfaceTypeValue.(map[string]interface{})


			if ok1 && i+1 < len(strValueSplit) {
			//	interfaceTypeValue = mapInterfaceTypeValue[strValueSplit[i+1]]
			//} else if ok2 {


				canonicalHeaderName:=strValueSplit[i+1]

				openingBracketIndex := strings.Index(canonicalHeaderName, "[")
				closingBracketIndex := strings.Index(canonicalHeaderName, "]")

				if openingBracketIndex == len(canonicalHeaderName)-3 && closingBracketIndex == len(canonicalHeaderName)-1 {
					// remove [2] part from Content-Type[2], set canonicalHeaderName to Content-Type
					canonicalHeaderName=canonicalHeaderName[0:openingBracketIndex]
					arrIndex := strValueSplit[i+1][openingBracketIndex+1:closingBracketIndex]
					arrIndexInt, err := strconv.Atoi(arrIndex) // convert string to int
					if err!=nil {
						log.Print("error parsing string to int ",err)
						return []string{}
					}
					var mapInterfaceArrayTypeValue []interface{}
					mapInterfaceArrayTypeValue, ok2 = (mapInterfaceTypeValue[canonicalHeaderName]).([]interface{})
					if ok2{
						interfaceTypeValue=mapInterfaceArrayTypeValue[arrIndexInt]
					}
				} else {
					// if the config is like $requestHeaders.Content-Type,return entire array of this header value
					interfaceTypeValue= mapInterfaceTypeValue[canonicalHeaderName]
				}


			} else {
				/* reached end of nested values i.e. requestJsonBody.[orderDetails].[addressDetails].pincode,
				so set the final value in response body map
				*/

				/*s,ok:=interfaceTypeValue.(string)
				if ok{
					return s
				}*/
				return interfaceTypeValue
			}
		}
	} else {
		return responseBodyConfigKeyValue
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
