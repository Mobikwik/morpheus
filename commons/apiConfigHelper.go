package commons

import (
	"encoding/json"
	"github.com/Mobikwik/morpheus/bboltDB"
	"github.com/Mobikwik/morpheus/model"
	"log"
)

func FindMatchingApiConfig(urlToSearch, requestMethod string) *model.ApiConfig {

	log.Printf("inside findMatchingApiConfig to find matching config for url %s requestMethod %s", urlToSearch, requestMethod)

	apiKey := makeApiConfigKey(urlToSearch, requestMethod)
	apiJsonFromDB := ReadSingleApiConfigFromDB(apiKey)
	if len(apiJsonFromDB) > 0 {
		var apiConfigJson model.ApiConfig
		json.Unmarshal([]byte(apiJsonFromDB), &apiConfigJson)
		log.Print("matching api config found with Id ", apiConfigJson.Id)
		return &apiConfigJson
	}
	log.Print("exiting findMatchingApiConfig")
	return nil
}

func makeApiConfigKey(urlToSearch, requestMethod string) string {
	return urlToSearch + "~" + requestMethod
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
