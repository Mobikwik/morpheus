package bboltDB

import (
	"encoding/json"
	"github.com/Mobikwik/morpheus/model"
	"go.etcd.io/bbolt"
	"log"
)

func UpdateApiConfigInDB(bucketName, key string, apiConfigObj model.ApiConfig) uint64 {

	log.Print("storing api config for key ", key)

	var id uint64
	var newApiConfigArr []model.ApiConfig
	err := boltDBConnection.Update(func(tx *bbolt.Tx) error {

		existingApiConfigJson, _ := ReadSingleKeyFromDB(bucketName, key)

		if len(existingApiConfigJson) > 0 {
			var existingApiConfigArr []model.ApiConfig
			json.Unmarshal([]byte(existingApiConfigJson), &existingApiConfigArr)
			arrLen := len(existingApiConfigArr)

			newApiConfigArr = make([]model.ApiConfig, arrLen+1)
			newApiConfigArr = existingApiConfigArr[0:arrLen]
			// add new value in arr
			newApiConfigArr[arrLen+1] = apiConfigObj

		} else {
			newApiConfigArr = []model.ApiConfig{apiConfigObj}
		}
		// bucket must be created/opened in same tx, hence passing tx in createBucket
		bucket := createBucket(bucketName, tx)
		id, _ = bucket.NextSequence()
		apiConfigObj.Id = id
		apiConfig, err := json.Marshal(&newApiConfigArr)
		if nil != err {
			log.Printf("error occured while updating DB %v", err)
			panic(err)
		}
		err = bucket.Put([]byte(key), apiConfig)

		if nil != err {
			log.Printf("error occured while updating DB %v", err)
			return err
		}
		return nil
	})

	if nil != err {
		panic(err)
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
