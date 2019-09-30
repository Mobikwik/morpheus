package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func newRouter() *mux.Router  {

	fmt.Println("entering newRouter method")

	r := mux.NewRouter()
	r.HandleFunc("/mock", handler).Methods("GET")

	staticFileDirectory := http.Dir("./assets/")

	staticFileHandler := http.StripPrefix("", http.FileServer(staticFileDirectory))

	r.PathPrefix("/").Handler(staticFileHandler).Methods("GET")

	fmt.Println("exiting newRouter method")

	return r
}

func main() {

	fmt.Println("entering mocking main")

	r := newRouter()

	http.ListenAndServe(":8080",r)

	fmt.Println("exiting mocking main")
}

func handler(w http.ResponseWriter, r *http.Request) {
	var msg string = "Welcome to Morpheus! Please use complete api url for mocking instead of "
	var url string = r.URL.Path

	fmt.Fprintf(w,msg+url)
}