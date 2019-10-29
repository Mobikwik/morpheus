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

	bboltDB.OpenDBConnection()

	portNumber := os.Args[1]
	r := webHandlers.NewRouter()
	err := http.ListenAndServe(":"+portNumber, r)

	if nil != err {
		log.Printf("error in running morpheus %v on port number %s", err, portNumber)
	}
	defer bboltDB.CloseDBConnection()
	log.Print("exiting mocking main")
}
