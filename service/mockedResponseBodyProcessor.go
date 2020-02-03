package service

import (
	"log"
	"strconv"
	"strings"
)

func setResponseBodyMap(responseBodyConfigJsonMap map[string]interface{}, requestBodyJsonMap map[string]interface{}) {

	for key, responseBodyConfigValueGenericType := range responseBodyConfigJsonMap {
		log.Printf("getting value for key %s response body config value %v of type %T", key, responseBodyConfigValueGenericType,
			responseBodyConfigValueGenericType)

		var responseBodyValueArr []interface{}

		switch responseBodyConfigValue := responseBodyConfigValueGenericType.(type) {

		case []string:
			for i, responseMockValuesValueSingle := range responseBodyConfigValue {
				responseBodyValueArr = append(responseBodyValueArr,
					GetResponseBodyValueFromRequestBody(responseMockValuesValueSingle, requestBodyJsonMap))
				log.Printf("adding array value %v on index %d for header %s ", responseBodyValueArr[i], i, key)
			}
			responseBodyConfigJsonMap[key] = responseBodyValueArr
		case []interface{}:
			responseBodyConfigJsonMap[key] = ProcessResponseMockValuesArrayType(responseBodyConfigValue,
				requestBodyJsonMap)
		case string:
			responseBodyConfigJsonMap[key] = GetResponseBodyValueFromRequestBody(responseBodyConfigValue, requestBodyJsonMap)

		case map[string]interface{}:
			// when the value is map type(nested json object), do recursive call
			setResponseBodyMap(responseBodyConfigValue, requestBodyJsonMap)
		default:
			log.Printf("no processing needed for key %s response body config %v type %T", key,
				responseBodyConfigValueGenericType,
				responseBodyConfigValueGenericType)
		}
	}
}

func ProcessResponseMockValuesArrayType(responseBodyConfigValue []interface{}, requestBodyJsonMap map[string]interface{}) (responseBodyValueArr []interface{}) {
	for i, responseBodyConfigKeyValueSingle := range responseBodyConfigValue {
		log.Printf("getting value for config %v", responseBodyConfigKeyValueSingle)
		responseBodyValueArr = append(responseBodyValueArr,
			GetResponseBodyValueFromRequestBody(responseBodyConfigKeyValueSingle, requestBodyJsonMap))
		log.Printf("adding array value %v on index %d for response body config %v", responseBodyValueArr[i], i,
			responseBodyConfigKeyValueSingle)

		/*	responseBodyConfigValueSingleStr, ok := responseBodyConfigKeyValueSingle.(string)
			if ok {
				responseBodyValueArr = append(responseBodyValueArr,
					GetResponseBodyValueFromRequestBody(responseBodyConfigValueSingleStr, requestBodyJsonMap))
				log.Printf("adding array value %v on index %d for response body config %s", responseBodyValueArr[i], i, responseBodyConfigValueSingleStr)
			} else {
				// TODO throw error
			}*/
	}
	return responseBodyValueArr
}

func GetResponseBodyValueFromRequestBody(responseBodyConfigValue interface{}, requestBodyJsonMap map[string]interface{}) interface{} {
	log.Printf("inside getResponseBodyValueFromRequestBody, getting value for response config %s", responseBodyConfigValue)

	switch responseBodyConfigValueTyped := responseBodyConfigValue.(type) {

	case string:
		if strings.HasPrefix(responseBodyConfigValueTyped, "requestBodyMockValues.") {

			responseBodyConfigValueSplit := strings.Split(responseBodyConfigValueTyped, ".")

			// temp variables to hold values fetched from requestBodyJsonMap.
			// Declaring here because these can't be declared inside below for loop as they need to hold values of previous iteration
			var requestBodyValueMapOfInterfaceType map[string]interface{}
			var requestBodyValueInterfaceType interface{}
			ok1 := false

			/*	get the value of first nesting level object reference
					for ex: if config is requestBodyMockValues.orderDetails.addressDetails.pincode, get value of requestBodyMockValues["orderDetails"]
					and store in requestBodyValueInterfaceType
				   Value can be of any type (string,number or another nested object, so storing in interface{} type)
			*/
			requestBodyValueInterfaceType = requestBodyJsonMap[responseBodyConfigValueSplit[1]]

			/* process all the nested object references from 2nd level onwards by looping around the array split with seperator ".",
			   i.e. iteration i=1: get value of requestBodyMockValues.[orderDetails].[addressDetails] from requestBodyValueInterfaceType and store in requestBodyValueInterfaceType
					iteration i=2: get value of requestBodyMockValues.[orderDetails].[addressDetails].pincode from requestBodyValueInterfaceType and store in requestBodyValueInterfaceType
			*/
			for i := 1; i < len(responseBodyConfigValueSplit); i++ {

				/* requestBodyValueMapOfInterfaceType is typecast of requestBodyValueInterfaceType from interface{} type to map[string]interface{}
				 i.e. value of requestBodyMockValues.[orderDetails].[addressDetails] is a nested json, so storing this value in
				requestBodyValueMapOfInterfaceType to extract further values (like requestBodyMockValues.[orderDetails].[addressDetails].pincode) from it in next iteration
				*/
				// Don't use ":=" for value assignment as it will redeclare requestBodyValueMapOfInterfaceType as a new local variable in each iteration
				requestBodyValueMapOfInterfaceType, ok1 = requestBodyValueInterfaceType.(map[string]interface{})

				// checking if we have more nested config values at (i+1)th level
				if ok1 && i+1 < len(responseBodyConfigValueSplit) {

					jsonKeyName := responseBodyConfigValueSplit[i+1]

					/* if config is like $requestBodyMockValues.txnTypes[2], get the array index part,i.e. 2
					responseBodyConfigValueSplit[1] = txnTypes[2], len(responseBodyConfigValueSplit[1]) =11,so:
					openingBracketIndex(index of [) = 8 = 11-3
					closingBracketIndex(index of ]) = 10 = 11-1
					arrIndex will have value 2
					*/
					openingBracketIndex := strings.Index(jsonKeyName, "[")
					closingBracketIndex := strings.Index(jsonKeyName, "]")

					if openingBracketIndex == len(jsonKeyName)-3 && closingBracketIndex == len(jsonKeyName)-1 {
						// get substring "2" from txnTypes[2]
						arrIndex := jsonKeyName[openingBracketIndex+1 : closingBracketIndex]
						// get substring "txnTypes" from txnTypes[2] and set in jsonKeyName
						jsonKeyName = jsonKeyName[0:openingBracketIndex]

						arrIndexInt, err := strconv.Atoi(arrIndex) // convert string "2" to int
						if err != nil {
							log.Printf("error parsing string to int for invalid config %s %v", responseBodyConfigValue, err)
							//panic(err)
							return nil // invalid config
						}
						var interfaceArrayTypeValue []interface{}
						ok2 := false
						/* As we are using config like txnTypes[2], value of txnTypes must be an array,
						so typecast to []interface{} */
						interfaceArrayTypeValue, ok2 = (requestBodyValueMapOfInterfaceType[jsonKeyName]).([]interface{})
						if ok2 {
							requestBodyValueInterfaceType = interfaceArrayTypeValue[arrIndexInt]
						} else {
							//throw error
							log.Printf("invalid config %s, trying to typecast non-array to array", responseBodyConfigValue)
							return nil
						}
					} else if openingBracketIndex == -1 && closingBracketIndex == -1 {
						// set value for next iteration if the config is like $requestHeaders.Content-Type i.e. does not have "[" and "]"
						requestBodyValueInterfaceType = requestBodyValueMapOfInterfaceType[jsonKeyName]
					}
				} else {
					/* reached end of nested values i.e. requestBodyMockValues.[orderDetails].[addressDetails].pincode, no more iterations possible,
					so return the final value from request body */
					return requestBodyValueInterfaceType
				}
			}
		} else {
			// response config does not start with "requestBodyMockValues.", so it's hard-coded value, so return that value as it is
			return responseBodyConfigValue
		}
	default:
		// not string value,so must be a constant, return as it is
		return responseBodyConfigValue
	}
	return nil
}
