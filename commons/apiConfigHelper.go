package commons

import (
	"encoding/json"
	"github.com/Mobikwik/morpheus/bboltDB"
	"github.com/Mobikwik/morpheus/model"
	"log"
	"reflect"
	"strings"
)

func FindMatchingMockConfig(urlToSearch string, requestMethod string,
	requestBodyMap map[string]interface{}) *model.MockConfig {

	log.Printf("inside findMatchingMockConfig to find matching config for url %s", urlToSearch)

	apiJsonFromDB := ReadSingleMockConfigFromDB(urlToSearch)
	if len(apiJsonFromDB) > 0 {
		//TODO fetch matching request config from array of MockConfig
		var mockConfigArray []model.MockConfig
		json.Unmarshal([]byte(apiJsonFromDB), &mockConfigArray)
		var mockConfig model.MockConfig
		for _, mockConfig = range mockConfigArray {
			//configRequestHeader := mockConfig.RequestMockValues.RequestHeadersMockValues
			configRequestMethod := mockConfig.Method
			strings.Compare("", "")
			configRequestBody := mockConfig.RequestMockValues.RequestBodyMockValues
			//check if request method and body values matches with one of the config values
			if configRequestMethod == requestMethod && reflect.DeepEqual(configRequestBody, requestBodyMap) {
				return &mockConfig
				/*	 if request body matches, check for request header values. Actual request might contain extra headers,
				so we check if headers present in config exists in actual request(requestHeader)
				 so requestHeader must contain configRequestHeader
				if isSubMap(configRequestHeader, requestHeader) {
					log.Printf("found matching api config with id %v value %v", mockConfig.Id, mockConfig)
					return &mockConfig
				}*/
			}
		}
	}
	return nil
}

func StoreMockConfigInDB(requestBodyJsonString, apiKey string) uint64 {
	var mockConfigJson model.MockConfig
	json.Unmarshal([]byte(requestBodyJsonString), &mockConfigJson)
	// set unique id
	return bboltDB.UpdateMockConfigInDB("mockConfig", apiKey, mockConfigJson)
}

func ReadSingleMockConfigFromDB(apiKey string) string {
	log.Print("reading api config for apikey ", apiKey)
	data, _ := bboltDB.ReadSingleKeyFromDB("mockConfig", apiKey)
	return data
}

func ReadEntireMockConfigFromDB() map[string]string {
	return bboltDB.ReadAllKeysFromDB("mockConfig")
}
