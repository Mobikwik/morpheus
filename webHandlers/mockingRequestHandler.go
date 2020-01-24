package webHandlers

import (
	"fmt"
	"github.com/Mobikwik/morpheus/commons"
	"github.com/Mobikwik/morpheus/service"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
)

func recoverConfigError(w http.ResponseWriter) {
	if r := recover(); r != nil {
		log.Print("panic occurred ", r)
		debug.PrintStack()
		fmt.Fprintf(w, "%v", r)
	}
}

func mockingRequestHandler(w http.ResponseWriter, r *http.Request) {

	// this function will be called in case of any "panic"
	defer recoverConfigError(w)

	log.Print("inside mockingRequestHandler")

	var url = r.URL.Path

	if strings.EqualFold(url, "/") {
		log.Print("no url provided for mocking")
		fmt.Fprintf(w, "%s", "Welcome to Morpheus, an api mocking framework by Mobikwik."+
			"Please use complete api url for mocking instead of /."+
			"For example: \"http://localhost:8080/api/customer/getOrders\"")
		return
	} else {
		bodyBytes := commons.ReadFromRequestBody(r.Body)

		log.Printf("api to mock is %s request method is %s", url, r.Method)

		var responseBody, responseHeaders, responseHttpCode = service.DoMocking(url, r.Method, bodyBytes, r.Header)

		log.Printf("final mocked response \n body: %s \n headers: %s \n http response code: %v",
			responseBody, responseHeaders, responseHttpCode)

		for headerName, headerValue := range responseHeaders {
			w.Header()[headerName] = headerValue
		}

		// set response http code as mentioned in api config
		if 200 != responseHttpCode {
			w.WriteHeader(responseHttpCode)
			log.Printf("set http response code to %v", responseHttpCode)
		}
		// send final api response
		fmt.Fprintf(w, "%v", responseBody)
	}
	log.Print("exiting mockingRequestHandler")
}
