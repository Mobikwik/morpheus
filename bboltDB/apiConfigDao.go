package bboltDB

import (
	"encoding/json"
	"github.com/Mobikwik/morpheus/model"
	"go.etcd.io/bbolt"
	"log"
)

/*
We store api config as a key-value pair in BBoltDB. Key is the api url(string) and value is array of configs([]model.ApiConfig).
There can be multiple configs for one api.
ex:
/api/p/wallet/testDebit1:

[{
	"Id": 9457,
	"Url": "/api/p/wallet/testDebit1",
	"Method": "POST",
	"ResponseDelayInSeconds": 11,
	"RequestConfig": {
		"RequestHeaders": {
			"X-DeviceId": "device123"
		},
		"RequestJsonBody": {
			"action": "debit",
			"module": "wallet",
			"txnDetails": {
				"amount": 456,
				"orderId": "cbsvdsd"
			}
		}
	},
	"ResponseConfig": {
		"HttpCode": 200,
		"ResponseHeaders": {
			"Checksum": "fdjfn",
			"X-DeviceId": "requestHeaders.X-DeviceId"
		},
		"ResponseJsonBody": {
			"actionDone": "requestJsonBody.action",
			"statusMsg": "Debit Success"
		}
	}
}, {
	"Id": 9458,
	"Url": "/api/p/wallet/testDebit1",
	"Method": "POST",
	"ResponseDelayInSeconds": 5,
	"RequestConfig": {
		"RequestHeaders": {
			"X-DeviceId": "device657"
		},
		"RequestJsonBody": {
			"action": "debit",
			"module": "upi",
			"txnDetails": {
				"amount": 234,
				"orderId": "fwiurwi"
			}
		}
	},
	"ResponseConfig": {
		"HttpCode": 200,
		"ResponseHeaders": {
			"Checksum": "fggg",
			"X-DeviceId": "requestHeaders.X-DeviceId"
		},
		"ResponseJsonBody": {
			"actionDone": "requestJsonBody.action",
			"statusMsg": "Debit Failed"
		}
	}
}]
This function adds a new api config if there is already a config present.
*/
func UpdateApiConfigInDB(bucketName, key string, apiConfigObj model.ApiConfig) uint64 {

	log.Print("storing api config for key ", key)

	var id uint64
	var newApiConfigArr []model.ApiConfig
	err1 := boltDBConnection.Update(func(tx *bbolt.Tx) error {
		// bucket must be created/opened in same tx, hence passing tx in createBucket
		bucket := createBucket(bucketName, tx)
		id, _ = bucket.NextSequence()
		apiConfigObj.Id = id

		// if there is already a config for this api, add this config in the existing config array
		existingApiConfigJson, _ := ReadSingleKeyFromDB(bucketName, key)
		if len(existingApiConfigJson) > 0 {
			var existingApiConfigArr []model.ApiConfig
			json.Unmarshal([]byte(existingApiConfigJson), &existingApiConfigArr)
			newApiConfigArr = append(existingApiConfigArr, apiConfigObj)
		} else {
			newApiConfigArr = []model.ApiConfig{apiConfigObj}
		}
		apiConfig, _ := json.Marshal(&newApiConfigArr)
		err2 := bucket.Put([]byte(key), apiConfig)
		if nil != err2 {
			log.Printf("error occured while updating DB %v", err2)
			return err2
		}
		return nil
	})

	if nil != err1 {
		panic(err1)
	}
	log.Printf("api config stored in db with id %v", id)
	return id
}

func ReadSingleKeyFromDB(bucketName, key string) (string, error) {

	var data string
	err := boltDBConnection.View(func(tx *bbolt.Tx) error {
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
	return data, nil
}

func ReadAllKeysFromDB(bucketName string) map[string]string {

	data := make(map[string]string)

	boltDBConnection.View(func(tx *bbolt.Tx) error {
		//bucket:= createBucket(bucketName,tx)
		bucket := tx.Bucket([]byte(bucketName))

		bucket.ForEach(func(k, v []byte) error {
			data[string(k)] = string(v)
			return nil
		})
		return nil
	})
	return data
}
