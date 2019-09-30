package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type VariableConfig struct {
	Name,DataType,Regex string
	MinLength,MaxLength int
}

func readVariableConfig() []VariableConfig {

	variableConfigJson:=`[

{
	"name": "orderId",
	"dataType": "string",
	"minLength": 4,
	"maxLength": 20,
	"regex": "[a-zA-Z0-9]{20}"
}

,
{
	"name": "bal",
	"dataType": "double",
	"minvalue": 0,
	"maxValue": 100000,
	"regex": "[0-9.]+"
}
,
{
	"name": "AuthToken",
	"dataType": "string",
	"fixLength": 20,
	"regex": "[a-zA-Z0-9~]{20}"
}
,
{
	"name": "ClientId",
	"dataType": "integer",
	"possibleValues": [3,9]
}
	]`

	fmt.Println(variableConfigJson)

	var variableConfig []VariableConfig
	json.Unmarshal([]byte(variableConfigJson), &variableConfig)
	fmt.Println("variableConfig values: ", variableConfig)

	return variableConfig

}

func variableConfigWebGetHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("inside variableConfigWebGetHandler")

	var variableConfigArr []VariableConfig

	variableConfigArr = readVariableConfig()

	fmt.Println("parsed json of variables is ",variableConfigArr)

	fmt.Fprintf(w,"%v",variableConfigArr)

	fmt.Println("exiting variableConfigWebGetHandler")
}

/*func main() {
	readVariableConfig()
}*/