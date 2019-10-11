package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func mockingRequestHandler(w http.ResponseWriter, r *http.Request) {

	log.Print("inside mockingRequestHandler")

	var url = r.URL.Path

	if strings.EqualFold(url,"/") {
		log.Print("no url provided for mocking")
		fmt.Fprintf(w,"%s", "Welcome to Morpheus, an api mocking framework by Mobikwik." +
			"Please use complete api url for mocking instead of /." +
			"For example: \"http://localhost:8080/api/customer/getOrders\"")
		return
	} else {
		var body =r.Body
		var bodyBytes []byte
		if body != nil {
			var err error
			bodyBytes, err = ioutil.ReadAll(body)
			if err != nil {
				panic(err)
			}
		}

		log.Printf("api to mock is %s request method is %s",url,r.Method)

		var responseBody, responseHeaders = doMocking(url,r.Method,bodyBytes,r.Header)

		log.Printf("final mocked response \n body: %s \n headers: %s",responseBody,responseHeaders)

		for headerName,headerValue:= range responseHeaders {
			w.Header()[headerName]=headerValue
		}
		// send final api response
		fmt.Fprintf(w,"%v",responseBody)
	}
	log.Print("exiting mockingRequestHandler")
}
