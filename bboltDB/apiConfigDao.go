package bboltDB

import (
	"encoding/json"
	"github.com/Mobikwik/morpheus/model"
	"go.etcd.io/bbolt"
	"log"
)

/*
We store api config as a key-value pair in BBoltDB. Key is the api url(string) and value is array of configs([]model.MockConfig).
There can be multiple configs for one api.
ex:
/api/p/wallet/testDebit1:

[{
	"Id": 9457,
	"Url": "/api/p/wallet/testDebit1",
	"Method": "POST",
	"ResponseDelayInSeconds": 11,
	"RequestMockValues": {
		"RequestHeadersMockValues": {
			"X-DeviceId": "device123"
		},
		"RequestBodyMockValues": {
			"action": "debit",
			"module": "wallet",
			"txnDetails": {
				"amount": 456,
				"orderId": "cbsvdsd"
			}
		}
	},
	"ResponseMockValues": {
		"HttpCode": 200,
		"ResponseHeadersMockValues": {
			"Checksum": "fdjfn",
			"X-DeviceId": "requestHeaders.X-DeviceId"
		},
		"ResponseBodyMockValues": {
			"actionDone": "requestJsonBody.action",
			"statusMsg": "Debit Success"
		}
	}
}, {
	"Id": 9458,
	"Url": "/api/p/wallet/testDebit1",
	"Method": "POST",
	"ResponseDelayInSeconds": 5,
	"RequestMockValues": {
		"RequestHeadersMockValues": {
			"X-DeviceId": "device657"
		},
		"RequestBodyMockValues": {
			"action": "debit",
			"module": "upi",
			"txnDetails": {
				"amount": 234,
				"orderId": "fwiurwi"
			}
		}
	},
	"ResponseMockValues": {
		"HttpCode": 200,
		"ResponseHeadersMockValues": {
			"Checksum": "fggg",
			"X-DeviceId": "requestHeaders.X-DeviceId"
		},
		"ResponseBodyMockValues": {
			"actionDone": "requestJsonBody.action",
			"statusMsg": "Debit Failed"
		}
	}
}]
This function adds a new api config if there is already a config present.
*/
func UpdateMockConfigInDB(bucketName, key string, mockConfigObj model.MockConfig) uint64 {

	log.Print("storing api config for key ", key)

	var id uint64
	var newMockConfigArr []model.MockConfig
	err1 := boltDBConnection.Update(func(tx *bbolt.Tx) error {
		// bucket must be created/opened in same tx, hence passing tx in createBucket
		bucket := createBucket(bucketName, tx)
		id, _ = bucket.NextSequence()
		mockConfigObj.Id = id

		// if there is already a config for this api, add this config in the existing config array
		existingMockConfigJson, _ := ReadSingleKeyFromDB(bucketName, key)
		if len(existingMockConfigJson) > 0 {
			var existingMockConfigArr []model.MockConfig
			json.Unmarshal([]byte(existingMockConfigJson), &existingMockConfigArr)
			newMockConfigArr = append(existingMockConfigArr, mockConfigObj)
		} else {
			newMockConfigArr = []model.MockConfig{mockConfigObj}
		}
		mockConfig, _ := json.Marshal(&newMockConfigArr)
		err2 := bucket.Put([]byte(key), mockConfig)
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
