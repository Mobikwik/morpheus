package commons

import (
	"encoding/json"
	"github.com/Mobikwik/morpheus/bboltDB"
	"github.com/Mobikwik/morpheus/model"
	"log"
)

func FindMatchingApiConfig(urlToSearch string) *model.ApiConfig {

	log.Printf("inside findMatchingApiConfig to find matching config for url %s", urlToSearch)

	apiKey := makeApiConfigKey(urlToSearch)
	apiJsonFromDB := ReadSingleApiConfigFromDB(apiKey)
	if len(apiJsonFromDB) > 0 {
		//TODO fetch matching request config from array of ApiConfig
		var apiConfigJson model.ApiConfig
		json.Unmarshal([]byte(apiJsonFromDB), &apiConfigJson)
		log.Print("matching api config found with Id ", apiConfigJson.Id)
		return &apiConfigJson
	}
	log.Print("exiting findMatchingApiConfig")
	return nil
}

func makeApiConfigKey(urlToSearch string) string {
	return urlToSearch
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
