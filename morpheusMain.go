package main

import (
	"github.com/Mobikwik/morpheus/bboltDB"
	"github.com/Mobikwik/morpheus/envConfig"
	"github.com/Mobikwik/morpheus/webHandlers"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	log.Print("entering mocking main")

	if len(os.Args) != 2 {
		panic("pass properties file path as command line args")
	}

	propertyFilePath := os.Args[1]

	p := envConfig.LoadProperties(propertyFilePath)

	dbPath := p.MustGetString("db.path") //"/tmp/bboltDBDataFile/morpheus.db"
	dbConnectTimeoutInSeconds := p.GetInt("db.connect.timeout", 1)
	bboltDB.OpenDBConnection(dbPath, dbConnectTimeoutInSeconds)

	portNumber := p.MustGetInt("webapp.port")
	r := webHandlers.NewRouter()
	err := http.ListenAndServe(":"+strconv.Itoa(portNumber), r)

	if nil != err {
		log.Printf("error in running morpheus %v on port number %s", err, portNumber)
	}
	defer bboltDB.CloseDBConnection()
	log.Print("exiting mocking main")
}
