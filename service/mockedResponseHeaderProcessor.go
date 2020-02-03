package service

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

func setResponseHeaderMap(responseHeaderConfigJsonMap map[string]interface{}, requestHeaderMap map[string][]string) (responseHeaderValuesJsonMap map[string][]string) {

	for headerName, responseMockValuesKeyValueGenericType := range responseHeaderConfigJsonMap {

		log.Printf("setting value for response header %s of type %T config value %s",
			headerName, responseMockValuesKeyValueGenericType, responseMockValuesKeyValueGenericType)

		var responseHeaderValueArr []string

		switch responseMockValuesKeyValue := responseMockValuesKeyValueGenericType.(type) {
		case string:
			responseHeaderValueArr = GetResponseHeaderConfigValueFromRequestHeader(responseMockValuesKeyValue, requestHeaderMap)
			log.Printf("setting single value %s for header %s", responseHeaderValueArr, headerName)
		/*case []string:
		for i, responseMockValuesKeyValueSingle := range responseMockValuesKeyValue {
			responseHeaderValueArr = append(responseHeaderValueArr, getResponseHeaderConfigValueFromRequestHeader(responseMockValuesKeyValueSingle, requestHeaderMap)[0])
			log.Printf("adding array value %s on index %d for header %s ", responseHeaderValueArr[i],i,headerName)
		}*/
		case []interface{}:
			log.Printf("handling response header value config for header %s config type is %T", headerName, responseMockValuesKeyValueGenericType)

			for i, responseMockValuesKeyValueSingle := range responseMockValuesKeyValue {
				log.Printf("getting value for config %s", responseMockValuesKeyValueSingle)
				responseMockValuesKeyValueSingleStr, ok := responseMockValuesKeyValueSingle.(string)
				if ok {
					responseHeaderValueArr = append(responseHeaderValueArr, GetResponseHeaderConfigValueFromRequestHeader(responseMockValuesKeyValueSingleStr, requestHeaderMap)[0])
					log.Printf("adding array value %s on index %d for header %s ", responseHeaderValueArr[i], i, headerName)
				}

			}

		default:
			log.Printf("invalid response header value config for header %s config type is %T", headerName, responseMockValuesKeyValueGenericType)

		}

		log.Printf("setting final value for response header %s = %s", headerName, responseHeaderValueArr)
		//responseHeaderConfigJsonMap[headerName] = responseHeaderValueArr
		responseHeaderValuesJsonMap[headerName] = responseHeaderValueArr
	}
	return responseHeaderValuesJsonMap
}

func GetResponseHeaderConfigValueFromRequestHeader(responseHeaderConfigValue string,
	requestHeaderMap map[string][]string) []string {
	log.Printf("inside getResponseHeaderConfigValueFromRequestHeader, getting value for response config %s",
		responseHeaderConfigValue)
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
			canonicalHeaderName = canonicalHeaderName[0:openingBracketIndex]
			arrIndex := responseHeaderConfigValueSplit[1][openingBracketIndex+1 : closingBracketIndex]
			arrIndexInt, err := strconv.Atoi(arrIndex) // convert string to int
			if err != nil {
				log.Print("error parsing string to int ", err)
				panic(err)
			}
			return []string{requestHeaderMap[canonicalHeaderName][arrIndexInt]}
			/*if arrIndexInt<len(requestHeaderMap[canonicalHeaderName]) {
				return []string{requestHeaderMap[canonicalHeaderName][arrIndexInt]}
			} else {
				log.Printf("invalid response header config %s",responseHeaderConfigValue)
				panic("invalid response header config "+responseHeaderConfigValue)
			}*/

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
