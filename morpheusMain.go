package main

import (
	"log"
	"github.com/gorilla/mux"
	"net/http"
)

func newRouter() *mux.Router  {

	log.Print("entering newRouter method")

	r := mux.NewRouter()

	// Handlers for APIs that return configs
	r.HandleFunc("/variableConfig", variableConfigWebGetHandler).Methods("GET")
	r.HandleFunc("/apiConfig", apiConfigWebGetHandler).Methods("GET")

	// Handler for mocking
	r.PathPrefix("/").HandlerFunc(mockingRequestHandler).Methods("GET","POST")

	staticFileDirectory := http.Dir("./assets/")

	staticFileHandler := http.StripPrefix("", http.FileServer(staticFileDirectory))

	r.PathPrefix("/").Handler(staticFileHandler).Methods("GET")

	log.Print("exiting newRouter method")

	return r
}

func main() {
	log.Print("entering mocking main")
	r := newRouter()
	http.ListenAndServe(":8080",r)
	log.Print("exiting mocking main")
}
