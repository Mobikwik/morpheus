package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Request struct{

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


func readApiConfigFromDB() string  {

	apiConfigJson:=`[
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
}
	]`

	return apiConfigJson

}

func getApiConfigArray() []ApiConfig  {

	var apiConfigJson = readApiConfigFromDB()

	return parseApiConfig(apiConfigJson)

}

func parseApiConfig(apiConfigJson string) []ApiConfig {

	log.Print(apiConfigJson)

	var apiConfig []ApiConfig
	json.Unmarshal([]byte(apiConfigJson), &apiConfig)


	for i, v := range apiConfig {
		log.Print("apiConfig values for i= ",i, v)
		//log.Printf("name %s minLength %d",v.Name,*(v.MinLength))
	}

	log.Print("apiConfig values: ", apiConfig)

	return apiConfig

}

func apiConfigWebGetHandler(w http.ResponseWriter, r *http.Request) {

	log.Print("inside apiConfigWebGetHandler")

	var apiConfigArr []ApiConfig

	apiConfigArr = getApiConfigArray()

	log.Print("parsed json of variables is ",apiConfigArr)

	fmt.Fprintf(w,"%v",apiConfigArr)

	log.Print("exiting apiConfigWebGetHandler")
}

/*func main() {
	parseVariableConfig()
}*/