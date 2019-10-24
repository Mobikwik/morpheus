package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Request struct {
	ResponseDelayInSeconds int
	// Header config values can be of type string or []string.Hence using generic interface{} type
	RequestHeaders map[string]interface{}
	// request body can have many types as string,numeric,array,another struct etc.Hence using generic interface{} type
	RequestJsonBody map[string]interface{}
}

type Response struct {
	HttpCode int
	// Header config values can be of type string or []string.Hence using generic interface{} type
	ResponseHeaders map[string]interface{}
	// response body can have many types as string,numeric,array,another struct etc.Hence using generic interface{} type
	ResponseJsonBody map[string]interface{}
}

type ApiConfig struct {
	Id             string
	Url            string
	Method         string
	RequestConfig  Request
	ResponseConfig Response
}

func readApiConfigFromDB() string {

	//TODO how to read entire config from DB
	testApiConfigJson := `[
{
	"id":"1",	
	"url": "/api/p/wallet/debit",
	"method": "POST",
	"requestConfig": {
		"responseDelayInSeconds": 30,
		"requestHeaders": {
			"Content-Type": ["application/json","text/html","application/pdf"],
			"Authorization": "$Auth",
			"X-DeviceId": "$DeviceId",
			"X-ClientId": "$ClientId",
			"X-Checksum": "hfsdhfbudgwq8gdqwudqu"
		},
		"requestJsonBody": {
			"action": "debit",
			"module": "wallet",
			"txnDetails": {
				"orderId": "$orderId",
				"amount": "$amt",
				"txnTypes":[1,2,3,4]
			}
		}
	},
	"responseConfig": {
		"httpCode": 200,
		"responseHeaders": {
			"X-DeviceId": "requestHeaders.X-DeviceId",
			"X-ClientId": "requestHeaders.X-ClientId",
			"Content-Type": "requestHeaders.Content-Type[0]",
			"AllContent-Types": "requestHeaders.Content-Type",
			"ConsideredContent-Types": ["requestHeaders.Content-Type[0]","requestHeaders.Content-Type[1]"],
			"SelectedContent-Type": "requestHeaders.Content-Type[0]",
			"DummyContent-Type": ["requestHeaders.Content-Type[0]","DummyContentTypeValue"],
			"Checksum": "fdjfnfffewfwef"
		},
		
		"responseJsonBody": {
			"actionDone": "requestJsonBody.action",
			"statusCode": "$statusCode",
			"statusMsg": "Debit Success",
			"orderId": "requestJsonBody.txnDetails.orderId",
			"consideredTxnTypes":["requestJsonBody.txnDetails.txnTypes[0]","requestJsonBody.txnDetails.txnTypes[1]"],
			"processedTxnType":	"requestJsonBody.txnDetails.txnTypes[0]",
			"allTxnTypes":"requestJsonBody.txnDetails.txnTypes",
			"amountDetails":{
				"amountDebited":"requestJsonBody.txnDetails.amount"
			},
			"requestTxnDetails":"requestJsonBody.txnDetails",
			"balanceData": {
				"mainBalance": 1023,
				"buckets": {
					"bucket1": "$bal",
					"bucket2": "$bal",
					"bucket3": "$bal"
				}
			}
		}
	}
},

{
	"id":"2",	
	"url": "/api/p/wallet/credit",
	"method": "POST",
	"requestConfig": {
		"responseDelayInSeconds": 30,
		"requestHeaders": {
			"Content-Type": ["application/json","text/html","application/pdf"],
			"Authorization": "$Auth",
			"X-DeviceId": "$DeviceId",
			"X-ClientId": "$ClientId",
			"X-Checksum": "hfsdhfbudgwq8gdqwudqu"
		},
		"requestJsonBody": {
			"action": "debit",
			"module": "wallet",
			"txnDetails": {
				"orderId": "$orderId",
				"amount": "$amt",
				"txnTypes":[1,2,3,4]
			}
		}
	},
	"responseConfig": {
		"httpCode": 200,
		"responseHeaders": {
			"X-DeviceId": "requestHeaders.X-DeviceId",
			"X-ClientId": "requestHeaders.X-ClientId",
			"Content-Type": "requestHeaders.Content-Type[0]",
			"AllContent-Types": "requestHeaders.Content-Type",
			"ConsideredContent-Types": ["requestHeaders.Content-Type[0]","requestHeaders.Content-Type[1]"],
			"SelectedContent-Type": "requestHeaders.Content-Type[0]",
			"DummyContent-Type": ["requestHeaders.Content-Type[0]","DummyContentTypeValue"],
			"Checksum": "fdjfnfffewfwef"
		},
		
		"responseJsonBody": {
			"actionDone": "requestJsonBody.action",
			"statusCode": "$statusCode",
			"statusMsg": "Debit Success",
			"orderId": "requestJsonBody.txnDetails.orderId",
			"consideredTxnTypes":["requestJsonBody.txnDetails.txnTypes[0]","requestJsonBody.txnDetails.txnTypes[1]"],
			"processedTxnType":	"requestJsonBody.txnDetails.txnTypes[0]",
			"allTxnTypes":"requestJsonBody.txnDetails.txnTypes",
			"amountDetails":{
				"amountDebited":"requestJsonBody.txnDetails.amount"
			},
			"requestTxnDetails":"requestJsonBody.txnDetails",
			"balanceData": {
				"mainBalance": 1023,
				"buckets": {
					"bucket1": "$bal",
					"bucket2": "$bal",
					"bucket3": "$bal"
				}
			}
		}
	}
}
	]`
	return testApiConfigJson
}

func getApiConfigArray() []ApiConfig {
	var apiConfigJson = readApiConfigFromDB()
	return parseApiConfig(apiConfigJson)
}

func parseApiConfig(apiConfigJson string) []ApiConfig {
	log.Print(apiConfigJson)
	var apiConfig []ApiConfig
	json.Unmarshal([]byte(apiConfigJson), &apiConfig)
	for i, v := range apiConfig {
		log.Print("apiConfig values for i= ", i, v)
	}
	log.Print("apiConfig values: ", apiConfig)
	return apiConfig
}

// Return API config stored in DB
func apiConfigWebGetHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("inside apiConfigWebGetHandler")
	queryStringValues := r.URL.Query()
	apiKey := ""
	if len(queryStringValues) >= 2 {
		apiKey = queryStringValues["apiUrl"][0] + "~" + queryStringValues["requestMethod"][0]
	}
	// fetch all configs
	if len(apiKey) == 0 {
		var apiConfigArr []ApiConfig
		apiConfigArr = getApiConfigArray()
		log.Print("parsed json of api config is ", apiConfigArr)
		fmt.Fprintf(w, "%v", apiConfigArr)
	} else {
		data := readApiConfigFromDB2(apiKey)
		fmt.Fprintf(w, "%v", data)
	}
	log.Print("exiting apiConfigWebGetHandler")
}

// Create new API config
func apiConfigWebPostHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("inside apiConfigWebPostHandler")

	requestBody := readFromRequestBody(r.Body)
	requestBodyJsonString := string(requestBody)

	var newApiConfig ApiConfig
	err := json.Unmarshal(requestBody, &newApiConfig)
	if err != nil {
		panic(err)
	}
	log.Println("parsed request body json for new api config is ", requestBodyJsonString)

	apiKey := newApiConfig.Url + "~" + newApiConfig.Method
	err = storeApiConfigInDB(requestBodyJsonString, apiKey)

	log.Print("exiting apiConfigWebPostHandler")
}

func storeApiConfigInDB(requestBodyJsonString, apiKey string) error {
	return updateInDB("mockApiConfig", apiKey, requestBodyJsonString)
}

func readApiConfigFromDB2(apiKey string) string {
	data, _ := read("mockApiConfig", apiKey)
	return data
}

func findMatchingApiConfig(urlToSearch, requestMethod string) *ApiConfig {
	//var matchingApiConfig *ApiConfig

	log.Printf("inside findMatchingApiConfig to find matching config for url %s requestMethod %s", urlToSearch, requestMethod)

	var apiConfigArr = getApiConfigArray()

	if apiConfigArr != nil {

		for _, apiConfig := range apiConfigArr {

			if strings.EqualFold(apiConfig.Url, urlToSearch) && strings.EqualFold(apiConfig.Method, requestMethod) {
				log.Print("matching api config found with Id ", apiConfig.Id)
				return &apiConfig
			}

		}
	}

	log.Print("exiting findMatchingApiConfig")

	return nil
}
