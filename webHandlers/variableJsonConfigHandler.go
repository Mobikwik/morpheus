package webHandlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type VariableConfig struct {
	Id, Name, DataType, Regex                           string
	MinLength, MaxLength, FixLength, MinValue, MaxValue *int
	PossibleStringValues                                []string
	PossibleIntValues                                   []int
}

func readVariableJsonFromDB() string {

	testVariableConfigJson := `[

{
	"id":"1",
	"name": "orderId",
	"dataType": "string",
	"minLength": 4,
	"maxLength": 20,
	"regex": "[a-zA-Z0-9]{20}"
}
,
{
	"id":"2",
	"name": "bal",
	"dataType": "double",
	"minValue": 0,
	"maxValue": 100000,
	"regex": "[0-9.]+"
}
,
{
	"id":"3",
	"name": "AuthToken",
	"dataType": "string",
	"fixLength": 20,
	"regex": "[a-zA-Z0-9~]{20}"
}
,
{
	"id":"4",
	"name": "ClientId",
	"dataType": "integer",
	"possibleIntValues": [3,9]
}
	]`

	return testVariableConfigJson
}

func parseVariableConfig(variableConfigJson string) []VariableConfig {

	log.Print(variableConfigJson)

	var variableConfig []VariableConfig
	json.Unmarshal([]byte(variableConfigJson), &variableConfig)

	for i, v := range variableConfig {
		log.Print("variableConfig values for i= ", i, v)
		//log.Printf("name %s minLength %d",v.Name,*(v.MinLength))
	}

	log.Print("variableConfig values: ", variableConfig)

	return variableConfig

}

func variableConfigWebGetHandler(w http.ResponseWriter, r *http.Request) {

	log.Print("inside variableConfigWebGetHandler")

	var variableConfigArr []VariableConfig

	var variableConfigJson = readVariableJsonFromDB()
	variableConfigArr = parseVariableConfig(variableConfigJson)

	log.Print("parsed json of variables is ", variableConfigArr)

	fmt.Fprintf(w, "%v", variableConfigArr)

	log.Print("exiting variableConfigWebGetHandler")
}

/*func main() {
	parseVariableConfig()
}*/
