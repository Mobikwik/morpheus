package main

import (
	"encoding/json"
	"github.com/Mobikwik/morpheus/commons"
	"github.com/Mobikwik/morpheus/service"
	"testing"
)

func TestResponseBodyConfig_ForSimpleValueFromRequestBody(t *testing.T) {
	responseBodyConfigValue := "simpleValue"
	expected := "simpleValue"

	runResponseBodyConfigTest(responseBodyConfigValue, expected, t)
}

func TestResponseBodyConfig_ForNestedValueFromRequestBody(t *testing.T) {

	responseBodyConfigValue := "requestJsonBody.txnDetails.orderId"
	expected := "abcd"

	runResponseBodyConfigTest(responseBodyConfigValue, expected, t)
}

func TestResponseBodyConfig_ForNestedSingleArrayValueFromRequestBody(t *testing.T) {

	responseBodyConfigValue := "requestJsonBody.txnDetails.txnTypes[2]"
	expected := "2"

	runResponseBodyConfigTest(responseBodyConfigValue, expected, t)
}

func TestResponseBodyConfig_ForNestedMultipleArrayValueFromRequestBody(t *testing.T) {

	responseBodyConfigValue := []interface{}{
		"requestJsonBody.txnDetails.txnTypes[2]",
		"requestJsonBody.txnDetails.txnTypes[3]",
	}
	expected := []interface{}{"2", "3"}

	runResponseBodyArrayTypeConfigTest(responseBodyConfigValue, expected, t)
}

func TestResponseBodyConfig_ForNestedEntireArrayFromRequestBody(t *testing.T) {

	responseBodyConfigValue := "requestJsonBody.txnDetails.txnTypes"
	expected := []interface{}{"0", "1", "2", "3"}

	runResponseBodyConfigTest(responseBodyConfigValue, expected, t)
}

func TestResponseBodyConfig_ForObjectValueFromRequestBody(t *testing.T) {

	responseBodyConfigValue := "requestJsonBody.txnDetails"
	expected := map[string]interface{}{
		"orderId":  "abcd",
		"amount":   "123",
		"txnTypes": []interface{}{"0", "1", "2", "3"},
	}
	runResponseBodyConfigTest(responseBodyConfigValue, expected, t)
}

func runResponseBodyArrayTypeConfigTest(responseBodyConfigValue []interface{},
	expected []interface{}, t *testing.T) {

	requestJson := `{ "action":"debit","txnDetails": {"orderId": "abcd","amount": 123,"txnTypes":["0","1","2","3"]}}`
	var requestBodyJsonMap map[string]interface{}
	err := json.Unmarshal([]byte(requestJson), &requestBodyJsonMap)
	if err != nil {
		t.Errorf("error in parsing json")
	}

	actual := service.processResponseConfigArrayType(responseBodyConfigValue, requestBodyJsonMap)

	if len(expected) != len(actual) {
		t.Errorf("expected array value %v type %T actual value %v type %T", expected, expected, actual, actual)
	}
	for i, v := range expected {
		if v != actual[i] {
			t.Errorf("expected array type value %v type %T actual value %v type %T", v, v, actual[i], actual[i])
		}
	}
}

func runResponseBodyConfigTest(responseBodyConfigValue string, expected interface{},
	t *testing.T) {

	requestJson := `{ "action":"debit","txnDetails": {"orderId": "abcd","amount": "123","txnTypes":["0","1","2","3"]}}`
	var requestBodyJsonMap map[string]interface{}
	err := json.Unmarshal([]byte(requestJson), &requestBodyJsonMap)
	if err != nil {
		t.Errorf("error in parsing json")
	}
	actual := service.getResponseBodyValueFromRequestBody(responseBodyConfigValue, requestBodyJsonMap)

	commons.compareValues(expected, actual, t)
}
