package bboltDB

import (
	"encoding/json"
	"github.com/Mobikwik/morpheus/model"
	"go.etcd.io/bbolt"
	"log"
)

func UpdateApiConfigInDB(bucketName, key string, apiConfig model.ApiConfig) {

	log.Print("storing api config for key ", key)

	err := boltDBConnection.Update(func(tx *bbolt.Tx) error {
		// bucket must be created/opened in same tx, hence passing tx in createBucket
		bucket := createBucket(bucketName, tx)
		id, _ := bucket.NextSequence()
		apiConfig.Id = id
		apiConfig, err := json.Marshal(&apiConfig)
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
