package main

import (
	bolt "go.etcd.io/bbolt"
	"log"
	"time"
)

func main() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("bboltDB/morpheus.db", 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		log.Fatal(err)
	} else {
		log.Print("BBoltDB connection opened successfully")
	}
	defer db.Close()
}
