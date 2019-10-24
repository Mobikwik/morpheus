package main

import (
	bolt "go.etcd.io/bbolt"
	"log"
	"time"
)

var db2 *bolt.DB

func createDBConnection() *bolt.DB {
	//TODO take DB file path, timeout etc from env.properties file
	db, err := bolt.Open("bboltDB/morpheus.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Print(err)
	} else {
		log.Printf("BBoltDB connection opened successfully to path %s", db.Path())
	}
	return db

}

func updateInDB(bucketName string, key, data string) error {

	log.Print("storing api config for key ", key)
	db := createDBConnection()

	err := db.Update(func(tx *bolt.Tx) error {
		bucket := createBucket(bucketName, tx)
		err := bucket.Put([]byte(key), []byte(data))

		if nil != err {
			log.Printf("error occured while updating DB %v", err)
			return err
		}
		return nil
	})

	if nil != err {
		return err
	}
	defer closeDBConnection(db)
	return nil
}

func read(bucketName, key string) (string, error) {

	/*if nil==db{
		createDBConnection()
		defer db.Close()
	}*/

	db := createDBConnection()

	var data string
	err := db.View(func(tx *bolt.Tx) error {
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
	defer closeDBConnection(db)
	return data, nil
}

func readAllKeysInMap(bucketName string) map[string]string {

	data := make(map[string]string)
	db := createDBConnection()

	db.View(func(tx *bolt.Tx) error {
		//bucket:= createBucket(bucketName,tx)
		bucket := tx.Bucket([]byte(bucketName))

		bucket.ForEach(func(k, v []byte) error {
			data[string(k)] = string(v)
			return nil
		})
		return nil
	})

	defer closeDBConnection(db)
	return data
}

func closeDBConnection(db *bolt.DB) {
	db.Close()
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
