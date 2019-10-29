package bboltDB

import (
	"encoding/json"
	"github.com/Mobikwik/morpheus/model"
	bolt "go.etcd.io/bbolt"
	"log"
)

func UpdateApiConfigInDB(bucketName, key string, apiConfig model.ApiConfig) {

	log.Print("storing api config for key ", key)

	readWriteDBConnection := createReadWriteDBConnection()

	err := readWriteDBConnection.Update(func(tx *bolt.Tx) error {
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
	defer CloseDBConnection(readWriteDBConnection)
}

func ReadAllKeysFromDB(bucketName string) map[string]string {

	data := make(map[string]string)
	readOnlyDBConnection := CreateReadOnlyDBConnection()

	readOnlyDBConnection.View(func(tx *bolt.Tx) error {
		//bucket:= createBucket(bucketName,tx)
		bucket := tx.Bucket([]byte(bucketName))

		bucket.ForEach(func(k, v []byte) error {
			data[string(k)] = string(v)
			return nil
		})
		return nil
	})
	defer CloseDBConnection(readOnlyDBConnection)
	return data
}
