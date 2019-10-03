package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Request struct{

	ResponseDelayInSeconds int
	RequestHeaders map[string]string
	RequestJsonBody interface{}
}

type Response struct {

	HttpCode int
	ResponseHeaders map[string]string
	ResponseJsonBody interface{}
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
			"Content-Type": "application/json",
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
				"amount": "$amt"
			}
		}
	},
	"responseConfig": {
		"httpCode": 200,
		"responseHeaders": {
			"X-DeviceId": "requestHeaders.X-DeviceId",
			"X-ClientId": "requestHeaders.X-ClientId",
			"Checksum": "fdjfnfffewfwef"
		},
		"responseJsonBody": {
			"statusCode": "$statusCode",
			"statusMsg": "Debit Success",
			"orderId": "requestJsonBody.txnDetails.orderId",
			"balanceData": {
				"mainBalance": "$bal",
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