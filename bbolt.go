package main

import (
	"encoding/json"
	bolt "go.etcd.io/bbolt"
	"log"
	"time"
)

func createReadOnlyDBConnection() *bolt.DB {

	//TODO take DB file path, timeout etc from env.properties file
	readOnlyDBConnection, err := bolt.Open("bboltDB/morpheus.db", 0600,
		&bolt.Options{Timeout: 1 * time.Second, ReadOnly: true})
	if err != nil {
		log.Print(err)
		panic(err)
	} else {
		log.Printf("BBoltDB read only connection opened successfully to path %s", readOnlyDBConnection.Path())
		return readOnlyDBConnection
	}
}

func createReadWriteDBConnection() *bolt.DB {

	//TODO take DB file path, timeout etc from env.properties file
	readWriteDBConnection, err := bolt.Open("bboltDB/morpheus.db", 0600,
		&bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Print(err)
		panic(err)
	} else {
		log.Printf("BBoltDB read-write connection opened successfully to path %s", readWriteDBConnection.Path())
		return readWriteDBConnection
	}
}

func updateApiConfigInDB(bucketName, key string, apiConfig ApiConfig) {

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
	defer closeDBConnection(readWriteDBConnection)
}

func read(bucketName, key string) (string, error) {

	readOnlyDBConnection := createReadOnlyDBConnection()

	var data string
	err := readOnlyDBConnection.View(func(tx *bolt.Tx) error {
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
	defer closeDBConnection(readOnlyDBConnection)
	return data, nil
}

func readAllKeysInMap(bucketName string) map[string]string {

	data := make(map[string]string)
	readOnlyDBConnection := createReadOnlyDBConnection()

	readOnlyDBConnection.View(func(tx *bolt.Tx) error {
		//bucket:= createBucket(bucketName,tx)
		bucket := tx.Bucket([]byte(bucketName))

		bucket.ForEach(func(k, v []byte) error {
			data[string(k)] = string(v)
			return nil
		})
		return nil
	})
	defer closeDBConnection(readOnlyDBConnection)
	return data
}

func closeDBConnection(db *bolt.DB) {
	db.Close()
	db = nil
	log.Print("closed db connection")
}

func createBucket(bucketName string, tx *bolt.Tx) *bolt.Bucket {
	var errInner error
	bucket, errInner := tx.CreateBucketIfNotExists([]byte(bucketName))
	if errInner != nil {
		panic(errInner)
	}
	return bucket
}
