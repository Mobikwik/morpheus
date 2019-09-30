package main

import (
	"encoding/json"
	"fmt"
	)

type VariableConfig struct {
	Name,DataType,Regex string
	MinLength,MaxLength int
}

func readMasterConfig(){

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
	fmt.Printf("variableConfig values: %s", variableConfig)

}

func main() {
	readMasterConfig()
}