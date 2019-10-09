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

	var msg = "Welcome to Morpheus! Please use complete api url for mocking instead of "
	var url = r.URL.Path

	if strings.EqualFold(url,"/") {
		log.Print("no url provided for mocking")
		fmt.Fprintf(w,msg+url)
		return
	} else {
		var body =r.Body
		var bodyBytes []byte
		if body != nil {
			bodyBytesTmp, err := ioutil.ReadAll(body)
			if err != nil {
				panic(err)
			}
			bodyBytes=bodyBytesTmp
		}

		// Restore the io.ReadCloser to its original state
		//body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))// Use the content
		//var bodyString = string(bodyBytes)

		log.Printf("api to mock is %s request method is %s",url,r.Method)


		var responseBody, responseContentType, responseHeaders = doMocking(url,r.Method,bodyBytes,r.Header)

		for k,v:= range responseHeaders {
			w.Header()[k]=v

		}
		log.Println("mocked response is ",responseBody,responseContentType,responseHeaders)
	}

	log.Print("exiting mockingRequestHandler")
}

/*func main() {
	parseVariableConfig()
}*/