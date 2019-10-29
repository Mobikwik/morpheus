package bboltDB

import (
	"go.etcd.io/bbolt"
	"log"
	"time"
)

var (
	readOnlyDBConnection *bbolt.DB
	readWriteDBConnection *bbolt.DB
	err error
)

func CreateReadOnlyDBConnection() *bbolt.DB {

	/*if nil!=readOnlyDBConnection{
		return readOnlyDBConnection
	}*/

	//TODO take DB file path, timeout etc from env.properties file
	readOnlyDBConnection, err := bbolt.Open("bboltDBDataFile/morpheus.db", 0600,
		&bbolt.Options{Timeout: 1 * time.Second, ReadOnly: true})
	if err != nil {
		log.Print(err)
		panic(err)
	} else {
		log.Printf("BBoltDB read only connection opened successfully to path %s", readOnlyDBConnection.Path())
		return readOnlyDBConnection
	}
}

func CloseDBConnection(db *bbolt.DB) {
	db.Close()
	db = nil
	log.Print("closed db connection")
}

func createBucket(bucketName string, tx *bbolt.Tx) *bbolt.Bucket {
	var errInner error
	bucket, errInner := tx.CreateBucketIfNotExists([]byte(bucketName))
	if errInner != nil {
		panic(errInner)
	}
	return bucket
}

func createReadWriteDBConnection() *bbolt.DB {

	/*if nil!=readWriteDBConnection {
		return readWriteDBConnection
	}*/

	//TODO take DB file path, timeout etc from env.properties file
	readWriteDBConnection, err = bbolt.Open("bboltDBDataFile/morpheus.db", 0600,
		&bbolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		log.Print(err)
		panic(err)
	} else {
		log.Printf("BBoltDB read-write connection opened successfully to path %s", readWriteDBConnection.Path())
		return readWriteDBConnection
	}
}
