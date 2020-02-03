package service_test

import (
	"encoding/json"
	"github.com/Mobikwik/morpheus/service"
	"reflect"
	"testing"
)

func TestResponseBodyConfig_ForSimpleValueFromRequestBody(t *testing.T) {
	responseBodyConfigValue := "simpleValue"
	expected := "simpleValue"
	runResponseBodyConfigTest(responseBodyConfigValue, expected, t)
}

func TestResponseBodyConfig_ForNestedValueFromRequestBody(t *testing.T) {

	responseBodyConfigValue := "requestBodyMockValues.txnDetails.orderId"
	expected := "abcd"

	runResponseBodyConfigTest(responseBodyConfigValue, expected, t)
}

func TestResponseBodyConfig_ForNestedSingleArrayValueFromRequestBody(t *testing.T) {

	responseBodyConfigValue := "requestBodyMockValues.txnDetails.txnTypes[2]"
	expected := "2"

	runResponseBodyConfigTest(responseBodyConfigValue, expected, t)
}

func TestResponseBodyConfig_ForNestedSingleArrayValueFromRequestBodyForInvalidConfigs(t *testing.T) {

	// Invalid array index
	responseBodyConfigValue := "requestBodyMockValues.txnDetails.txnTypes[a]"
	var expected interface{} = nil
	runResponseBodyConfigTest(responseBodyConfigValue, expected, t)

	// Invalid nested config
	responseBodyConfigValue = "requestBodyMockValues."
	expected = nil
	runResponseBodyConfigTest(responseBodyConfigValue, expected, t)
}

func TestResponseBodyConfig_ForNestedMultipleArrayValueFromRequestBody(t *testing.T) {

	responseBodyConfigValue := []interface{}{
		"requestBodyMockValues.txnDetails.txnTypes[2]",
		"requestBodyMockValues.txnDetails.txnTypes[3]",
	}
	expected := []interface{}{"2", "3"}
	runResponseBodyArrayTypeConfigTest(responseBodyConfigValue, expected, t)

	var parsedJson []interface{}

	jsonS := `["a", "2", "3c"]`
	json.Unmarshal([]byte(jsonS), &parsedJson)
	runResponseBodyArrayTypeConfigTest(parsedJson, parsedJson, t)

	jsonS = `[1, 2, 3]`
	json.Unmarshal([]byte(jsonS), &parsedJson)
	runResponseBodyArrayTypeConfigTest(parsedJson, parsedJson, t)

	jsonS = `["1", 2, "3c"]`
	json.Unmarshal([]byte(jsonS), &parsedJson)
	runResponseBodyArrayTypeConfigTest(parsedJson, parsedJson, t)

	jsonS = ` [

		{
			"orderId": "abc123",
			"amount": 23
		},
		{
			"orderId": "xyz456",
			"amount": 45
		},
		{
			"orderId": "jdfj4546",
			"amount": 789
		}

	] `

	json.Unmarshal([]byte(jsonS), &parsedJson)
	runResponseBodyArrayTypeConfigTest(parsedJson, parsedJson, t)
}

func TestResponseBodyConfig_ForNestedEntireArrayFromRequestBody(t *testing.T) {

	responseBodyConfigValue := "requestBodyMockValues.txnDetails.txnTypes"
	expected := []interface{}{"0", "1", "2", "3"}

	runResponseBodyConfigTest(responseBodyConfigValue, expected, t)
}

func TestResponseBodyConfig_ForObjectValueFromRequestBody(t *testing.T) {

	responseBodyConfigValue := "requestBodyMockValues.txnDetails"
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

	actual := service.ProcessResponseMockValuesArrayType(responseBodyConfigValue, requestBodyJsonMap)

	isEqual := reflect.DeepEqual(expected, actual)
	if !isEqual {
		t.Errorf("expected value %v type %T actual value %v type %T", expected, expected,
			actual, actual)
	}

	/*	if len(expected) != len(actual) {
			t.Errorf("expected array value %v type %T actual value %v type %T", expected, expected, actual, actual)
		}
		for i, v := range expected {
			if v != actual[i] {
				t.Errorf("expected array type value %v type %T actual value %v type %T", v, v, actual[i], actual[i])
			}
		}*/
}

func runResponseBodyConfigTest(responseBodyConfigValue string, expected interface{},
	t *testing.T) {

	requestJson := `{ "action":"debit","txnDetails": {"orderId": "abcd","amount": "123","txnTypes":["0","1","2","3"]}}`
	var requestBodyJsonMap map[string]interface{}
	err := json.Unmarshal([]byte(requestJson), &requestBodyJsonMap)
	if err != nil {
		t.Errorf("error in parsing json")
	}
	actual := service.GetResponseBodyValueFromRequestBody(responseBodyConfigValue, requestBodyJsonMap)

	//commons.CompareValues(expected, actual, t)
	isEqual := reflect.DeepEqual(expected, actual)
	if !isEqual {
		t.Errorf("expected value %v type %T actual value %v type %T", expected, expected,
			actual, actual)
	}
}
