package webHandlers

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func NewRouter() *mux.Router {

	log.Print("entering newRouter method")

	r := mux.NewRouter()

	// Handlers for APIs that return configs
	r.HandleFunc("/variableConfig", variableConfigWebGetHandler).Methods("GET")
	r.HandleFunc("/mockConfig", mockConfigWebGetHandler).Methods("GET")
	r.HandleFunc("/mockConfig", mockConfigWebPostHandler).Methods("POST")

	// Handler for mocking
	r.PathPrefix("/").HandlerFunc(mockingRequestHandler).Methods("GET", "POST")

	staticFileDirectory := http.Dir("./assets/")

	staticFileHandler := http.StripPrefix("", http.FileServer(staticFileDirectory))

	r.PathPrefix("/").Handler(staticFileHandler).Methods("GET")

	log.Print("exiting newRouter method")

	return r
}
