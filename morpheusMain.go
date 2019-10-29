package main

import (
	"github.com/Mobikwik/morpheus/webHandlers"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Print("entering mocking main")

	portNumber := os.Args[1]
	r := webHandlers.NewRouter()
	err := http.ListenAndServe(":"+portNumber, r)
	if nil != err {
		log.Printf("error in running morpheus %v on port number %s", err, portNumber)
	}
	log.Print("exiting mocking main")
}
