package main

import (
	"github.com/Mobikwik/morpheus/bboltDB"
	"github.com/Mobikwik/morpheus/webHandlers"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Print("entering mocking main")

	if len(os.Args) != 3 {
		panic("pass port number and db path as command line args")
	}

	portNumber := os.Args[1]
	//TODO take DB file path, timeout etc from env.properties file
	dbPath := os.Args[2] //"/tmp/bboltDBDataFile/morpheus.db"
	dbConnectTimeoutInSeconds := 1

	bboltDB.OpenDBConnection(dbPath, dbConnectTimeoutInSeconds)

	r := webHandlers.NewRouter()
	err := http.ListenAndServe(":"+portNumber, r)

	if nil != err {
		log.Printf("error in running morpheus %v on port number %s", err, portNumber)
	}
	defer bboltDB.CloseDBConnection()
	log.Print("exiting mocking main")
}
