package bboltDB

import (
	"go.etcd.io/bbolt"
	"log"
	"time"
)

var (
	boltDBConnection *bbolt.DB
	err              error
)

func OpenDBConnection(dbPath string, dbConnectTimeoutInSeconds int) {

	if nil != boltDBConnection {
		return
	}
	boltDBConnection, err = bbolt.Open(dbPath, 0600,
		&bbolt.Options{Timeout: time.Duration(dbConnectTimeoutInSeconds) * time.Second})
	if err != nil {
		log.Print(err)
		panic(err)
	} else {
		log.Printf("BBoltDB read-write connection opened successfully to path %s", boltDBConnection.Path())
	}
}

func createBucket(bucketName string, tx *bbolt.Tx) *bbolt.Bucket {
	var errInner error
	bucket, errInner := tx.CreateBucketIfNotExists([]byte(bucketName))
	if errInner != nil {
		panic(errInner)
	}
	return bucket
}

func CloseDBConnection() {
	boltDBConnection.Close()
	log.Print("closed db connection")
}