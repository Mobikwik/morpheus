package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func newRouter() *mux.Router {

	log.Print("entering newRouter method")

	r := mux.NewRouter()

	// Handlers for APIs that return configs
	r.HandleFunc("/variableConfig", variableConfigWebGetHandler).Methods("GET")
	r.HandleFunc("/apiConfig", apiConfigWebGetHandler).Methods("GET")
	r.HandleFunc("/apiConfig", apiConfigWebPostHandler).Methods("POST")

	// Handler for mocking
	r.PathPrefix("/").HandlerFunc(mockingRequestHandler).Methods("GET", "POST")

	staticFileDirectory := http.Dir("./assets/")

	staticFileHandler := http.StripPrefix("", http.FileServer(staticFileDirectory))

	r.PathPrefix("/").Handler(staticFileHandler).Methods("GET")

	log.Print("exiting newRouter method")

	return r
}

func main() {
	log.Print("entering mocking main")

	portNumber := os.Args[1]
	r := newRouter()
	err := http.ListenAndServe(":"+portNumber, r)
	if nil != err {
		log.Printf("error in running morpheus %v on port number %s", err, portNumber)
	}
	log.Print("exiting mocking main")
}
