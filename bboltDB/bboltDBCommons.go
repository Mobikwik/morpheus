package bboltDB

import (
	"go.etcd.io/bbolt"
	"log"
	"time"
)

var (
	boltDBConnection *bbolt.DB
)

func OpenDBConnection(dbPath string, dbConnectTimeoutInSeconds int) {

	if nil != boltDBConnection {
		return
	}
	var err error
	boltDBConnection, err = bbolt.Open(dbPath, 0600,
		&bbolt.Options{Timeout: time.Duration(dbConnectTimeoutInSeconds) * time.Second})
	if err != nil {
		log.Print(err)
		panic(err)
	} else {
		log.Printf("BBoltDB read-write connection opened successfully to path %s", boltDBConnection.Path())
	}
}

func createBucket(bucketName string, tx *bbolt.Tx) (bucket *bbolt.Bucket) {
	bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
	if err != nil {
		panic(err)
	}
	return bucket
}

func CloseDBConnection() {
	boltDBConnection.Close()
	log.Print("closed db connection")
}
