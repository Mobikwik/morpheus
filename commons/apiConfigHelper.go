package commons

import (
	"encoding/json"
	"github.com/Mobikwik/morpheus/bboltDB"
	"github.com/Mobikwik/morpheus/model"
	"log"
	"reflect"
)

func FindMatchingApiConfig(urlToSearch string, requestHeader map[string][]string,
	requestBodyMap map[string]interface{}) *model.ApiConfig {

	log.Printf("inside findMatchingApiConfig to find matching config for url %s", urlToSearch)

	apiJsonFromDB := ReadSingleApiConfigFromDB(urlToSearch)
	if len(apiJsonFromDB) > 0 {
		//TODO fetch matching request config from array of ApiConfig
		var apiConfigArray []model.ApiConfig
		json.Unmarshal([]byte(apiJsonFromDB), &apiConfigArray)
		var apiConfig model.ApiConfig
		for _, apiConfig = range apiConfigArray {
			//configRequestHeader := apiConfig.RequestMockValues.RequestHeadersMockValues
			configRequestBody := apiConfig.RequestMockValues.RequestBodyMockValues
			//check if request body values matches with one of the config values
			if reflect.DeepEqual(configRequestBody, requestBodyMap) {
				return &apiConfig
				/*	 if request body matches, check for request header values. Actual request might contain extra headers,
				so we check if headers present in config exists in actual request(requestHeader)
				 so requestHeader must contain configRequestHeader
				if isSubMap(configRequestHeader, requestHeader) {
					log.Printf("found matching api config with id %v value %v", apiConfig.Id, apiConfig)
					return &apiConfig
				}*/
			}
		}
	}
	return nil
}

func StoreApiConfigInDB(requestBodyJsonString, apiKey string) uint64 {
	var apiConfigJson model.ApiConfig
	json.Unmarshal([]byte(requestBodyJsonString), &apiConfigJson)
	// set unique id
	return bboltDB.UpdateApiConfigInDB("mockApiConfig", apiKey, apiConfigJson)
}

func ReadSingleApiConfigFromDB(apiKey string) string {
	log.Print("reading api config for apikey ", apiKey)
	data, _ := bboltDB.ReadSingleKeyFromDB("mockApiConfig", apiKey)
	return data
}

func ReadEntireApiConfigFromDB() map[string]string {
	return bboltDB.ReadAllKeysFromDB("mockApiConfig")
}
