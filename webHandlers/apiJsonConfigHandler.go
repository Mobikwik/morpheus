package webHandlers

import (
	"encoding/json"
	"fmt"
	"github.com/Mobikwik/morpheus/commons"
	"github.com/Mobikwik/morpheus/model"
	"log"
	"net/http"
)

/*func readMockConfigFromDB() string {

	testMockConfigJson := `[
{
	"id":"1",
	"url": "/api/p/wallet/debit",
	"method": "POST",
	"responseDelayInSeconds": 30,
	"requestMockValues": {
		"requestHeaders": {
			"Content-Type": ["application/json","text/html","application/pdf"],
			"Authorization": "$Auth",
			"X-DeviceId": "$DeviceId",
			"X-ClientId": "$ClientId",
			"X-Checksum": "hfsdhfbudgwq8gdqwudqu"
		},
		"requestBodyMockValues": {
			"action": "debit",
			"module": "wallet",
			"txnDetails": {
				"orderId": "$orderId",
				"amount": "$amt",
				"txnTypes":[1,2,3,4]
			}
		}
	},
	"responseMockValues": {
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
			"actionDone": "requestBodyMockValues.action",
			"statusCode": "$statusCode",
			"statusMsg": "Debit Success",
			"orderId": "requestBodyMockValues.txnDetails.orderId",
			"consideredTxnTypes":["requestBodyMockValues.txnDetails.txnTypes[0]","requestBodyMockValues.txnDetails.txnTypes[1]"],
			"processedTxnType":	"requestBodyMockValues.txnDetails.txnTypes[0]",
			"allTxnTypes":"requestBodyMockValues.txnDetails.txnTypes",
			"amountDetails":{
				"amountDebited":"requestBodyMockValues.txnDetails.amount"
			},
			"requestTxnDetails":"requestBodyMockValues.txnDetails",
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
	"requestMockValues": {
		"requestHeaders": {
			"Content-Type": ["application/json","text/html","application/pdf"],
			"Authorization": "$Auth",
			"X-DeviceId": "$DeviceId",
			"X-ClientId": "$ClientId",
			"X-Checksum": "hfsdhfbudgwq8gdqwudqu"
		},
		"requestBodyMockValues": {
			"action": "debit",
			"module": "wallet",
			"txnDetails": {
				"orderId": "$orderId",
				"amount": "$amt",
				"txnTypes":[1,2,3,4]
			}
		}
	},
	"responseMockValues": {
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
			"actionDone": "requestBodyMockValues.action",
			"statusCode": "$statusCode",
			"statusMsg": "Debit Success",
			"orderId": "requestBodyMockValues.txnDetails.orderId",
			"consideredTxnTypes":["requestBodyMockValues.txnDetails.txnTypes[0]","requestBodyMockValues.txnDetails.txnTypes[1]"],
			"processedTxnType":	"requestBodyMockValues.txnDetails.txnTypes[0]",
			"allTxnTypes":"requestBodyMockValues.txnDetails.txnTypes",
			"amountDetails":{
				"amountDebited":"requestBodyMockValues.txnDetails.amount"
			},
			"requestTxnDetails":"requestBodyMockValues.txnDetails",
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
	return testMockConfigJson
}
*/
/*
func getMockConfigArray() []MockConfig {
	var mockConfigJson = readMockConfigFromDB()
	return parseMockConfig(mockConfigJson)
}

func parseMockConfig(mockConfigJson string) []MockConfig {
	log.Print(mockConfigJson)
	var mockConfig []MockConfig
	json.Unmarshal([]byte(mockConfigJson), &mockConfig)
	for i, v := range mockConfig {
		log.Print("mockConfig values for i= ", i, v)
	}
	log.Print("mockConfig values: ", mockConfig)
	return mockConfig
}*/

// Return API config stored in DB
func mockConfigWebGetHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("inside mockConfigWebGetHandler")
	queryStringValues := r.URL.Query()
	var apiKey string
	if len(queryStringValues) >= 1 {
		// get values of query string params
		apiKey = queryStringValues["apiUrl"][0]
	}
	// fetch all configs
	if len(apiKey) == 0 {
		data := commons.ReadEntireMockConfigFromDB()
		fmt.Fprintf(w, "%v", data)
	} else {
		data := commons.ReadSingleMockConfigFromDB(apiKey)
		fmt.Fprintf(w, "%v", data)
	}
	log.Print("exiting mockConfigWebGetHandler")
}

// Create new API config
func mockConfigWebPostHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("inside mockConfigWebPostHandler")

	requestBody := commons.ReadFromRequestBody(r.Body)
	requestBodyJsonString := string(requestBody)

	var newMockConfig model.MockConfig
	err := json.Unmarshal(requestBody, &newMockConfig)
	if err != nil {
		panic(err)
	}
	log.Println("parsed request body json for new api config is ", requestBodyJsonString)

	apiKey := newMockConfig.Url
	id := commons.StoreMockConfigInDB(requestBodyJsonString, apiKey)
	fmt.Fprintf(w, "%v %d", "api config stored in DB with id ", id)

	log.Print("exiting mockConfigWebPostHandler")
}
