package commons

import (
	"encoding/json"
	"github.com/Mobikwik/morpheus/bboltDB"
	"github.com/Mobikwik/morpheus/model"
	"go.etcd.io/bbolt"
	"log"
)

func FindMatchingApiConfig(urlToSearch, requestMethod string) *model.ApiConfig {
	//var matchingApiConfig *ApiConfig

	log.Printf("inside findMatchingApiConfig to find matching config for url %s requestMethod %s", urlToSearch, requestMethod)

	apiKey := makeApiConfigKey(urlToSearch, requestMethod)
	apiJsonFromDB := ReadSingleApiConfigFromDB(apiKey)
	if len(apiJsonFromDB) > 0 {
		var apiConfigJson model.ApiConfig
		json.Unmarshal([]byte(apiJsonFromDB), &apiConfigJson)
		log.Print("matching api config found with Id ", apiConfigJson.Id)
		return &apiConfigJson
	}
	/*	var apiConfigArr = getApiConfigArray()
		if apiConfigArr != nil {
			for _, apiConfig := range apiConfigArr {
				if strings.EqualFold(apiConfig.Url, urlToSearch) && strings.EqualFold(apiConfig.Method, requestMethod) {
					log.Print("matching api config found with Id ", apiConfig.Id)
					return &apiConfig
				}
			}
		}*/

	log.Print("exiting findMatchingApiConfig")

	return nil
}

func makeApiConfigKey(urlToSearch, requestMethod string) string {
	return urlToSearch + "~" + requestMethod
}

func ReadSingleKeyFromDB(bucketName, key string) (string, error) {

	readOnlyDBConnection := bboltDB.CreateReadOnlyDBConnection()

	var data string
	err := readOnlyDBConnection.View(func(tx *bbolt.Tx) error {
		//bucket:= createBucket(bucketName,tx)
		bucket := tx.Bucket([]byte(bucketName))
		dataBytes := bucket.Get([]byte(key))
		data = string(dataBytes)
		return nil
	})

	if nil != err {
		log.Printf("error occured while reading from DB %v", err)
		return "", err
	}
	defer bboltDB.CloseDBConnection(readOnlyDBConnection)
	return data, nil
}

func StoreApiConfigInDB(requestBodyJsonString, apiKey string) {
	var apiConfigJson model.ApiConfig
	json.Unmarshal([]byte(requestBodyJsonString), &apiConfigJson)
	// set unique id
	bboltDB.UpdateApiConfigInDB("mockApiConfig", apiKey, apiConfigJson)
}

func ReadSingleApiConfigFromDB(apiKey string) string {
	data, _ := ReadSingleKeyFromDB("mockApiConfig", apiKey)
	return data
}

func ReadEntireApiConfigFromDB() map[string]string {
	return bboltDB.ReadAllKeysFromDB("mockApiConfig")
}
