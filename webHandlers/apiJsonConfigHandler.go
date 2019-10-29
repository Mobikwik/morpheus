package webHandlers

import (
	"encoding/json"
	"fmt"
	"github.com/Mobikwik/morpheus/commons"
	"github.com/Mobikwik/morpheus/model"
	"log"
	"net/http"
)

/*func readApiConfigFromDB() string {

	testApiConfigJson := `[
{
	"id":"1",
	"url": "/api/p/wallet/debit",
	"method": "POST",
	"responseDelayInSeconds": 30,
	"requestConfig": {
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
	"responseDelayInSeconds": 30,
	"requestConfig": {
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
*/
/*
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
}*/

// Return API config stored in DB
func apiConfigWebGetHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("inside apiConfigWebGetHandler")
	queryStringValues := r.URL.Query()
	apiKey := ""
	if len(queryStringValues) >= 2 {
		// get values of query string params
		apiKey = queryStringValues["apiUrl"][0] + "~" + queryStringValues["requestMethod"][0]
	}
	// fetch all configs
	if len(apiKey) == 0 {
		data := commons.ReadEntireApiConfigFromDB()
		fmt.Fprintf(w, "%v", data)
	} else {
		data := commons.ReadSingleApiConfigFromDB(apiKey)
		fmt.Fprintf(w, "%v", data)
	}
	log.Print("exiting apiConfigWebGetHandler")
}

// Create new API config
func apiConfigWebPostHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("inside apiConfigWebPostHandler")

	requestBody := commons.ReadFromRequestBody(r.Body)
	requestBodyJsonString := string(requestBody)

	var newApiConfig model.ApiConfig
	err := json.Unmarshal(requestBody, &newApiConfig)
	if err != nil {
		panic(err)
	}
	log.Println("parsed request body json for new api config is ", requestBodyJsonString)

	apiKey := newApiConfig.Url + "~" + newApiConfig.Method
	commons.StoreApiConfigInDB(requestBodyJsonString, apiKey)

	log.Print("exiting apiConfigWebPostHandler")
}

