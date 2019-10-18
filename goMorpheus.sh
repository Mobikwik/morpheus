export PATH=$PATH:/home/chiragjain/mobikwik/softwares/go1.12.5.linux-amd64/go/bin
export GOPATH=/home/chiragjain/mobikwik/codebase/goLangProjects
export CGO_ENABLED=0
go run morpheusMain.go mockedResponseBodyProcessor.go apiJsonConfigHandler.go mockedResponseProcessor.go variableJsonConfigHandler.go mockingRequestHandler.go mockedResponseHeaderProcessor.go commonTestFunctions.go 8080 &
