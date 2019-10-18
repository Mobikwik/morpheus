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
	log.Print("port to run morpheus is " + portNumber)

	r := newRouter()
	http.ListenAndServe(":"+portNumber, r)
	log.Print("exiting mocking main")
}
