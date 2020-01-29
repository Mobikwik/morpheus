package service_test

import (
	"github.com/Mobikwik/morpheus/service"
	"reflect"
	"testing"
)

func TestResponseHeaderConfig_ForSimpleValueFromRequestHeader(t *testing.T) {
	responseHeaderConfigValue := "simpleValue"
	expected := []string{"simpleValue"}

	runResponseHeaderConfigTest(responseHeaderConfigValue, expected, t)
}

func TestResponseHeaderConfig_ForNestedValueFromRequestHeader(t *testing.T) {

	responseHeaderConfigValue := "requestHeaders.X-DeviceId"
	expected := []string{"Device1234"}

	runResponseHeaderConfigTest(responseHeaderConfigValue, expected, t)
}

func TestResponseHeaderConfig_ForNestedSingleArrayValueFromRequestHeader(t *testing.T) {

	responseHeaderConfigValue := "requestHeaders.Content-Type[2]"
	expected := []string{"application/pdf"}

	runResponseHeaderConfigTest(responseHeaderConfigValue, expected, t)
}

func TestResponseHeaderConfig_ForNestedEntireArrayFromRequestHeader(t *testing.T) {

	responseHeaderConfigValue := "requestHeaders.Content-Type"
	expected := []string{"application/json", "text/html", "application/pdf"}

	runResponseHeaderConfigTest(responseHeaderConfigValue, expected, t)
}

// TODO test pending for array type header config
/*
func TestResponseHeaderConfig_ForNestedMultipleArrayValueFromRequestHeader(t *testing.T)  {
	responseHeaderConfigValue:= []interface{} {
		"requestHeaders.Content-Type[0]",
		"requestHeaders.Content-Type[1]",
	}

	expected := []string {"application/json","text/html"}

	runResponseHeaderArrayTypeConfigTest(responseHeaderConfigValue, expected,t)
	//runResponseHeaderConfigTest(responseHeaderConfigValue, expected,t)
}


func runResponseHeaderArrayTypeConfigTest(responseHeaderConfigValue []interface{},
	expected []string, t *testing.T) {

	requestJson:=`{ "action":"debit","txnDetails": {"orderId": "abcd","amount": 123,"txnTypes":["0","1","2","3"]}}`
	var requestHeaderJsonMap map[string]interface{}
	err := json.Unmarshal([]byte(requestJson), &requestHeaderJsonMap)
	if err != nil {
		t.Errorf("error in parsing json")
	}

	actual := processResponseMockValuesArrayType(responseHeaderConfigValue, requestHeaderJsonMap)

	if len(expected)!=len(actual){
		t.Errorf("expected array value %v type %T actual value %v type %T", expected,expected, actual, actual)
	}
	for i,v:=range expected{
		if v!= actual[i] {
			t.Errorf("expected array type value %v type %T actual value %v type %T", v,v, actual[i], actual[i])
		}
	}
}*/

func runResponseHeaderConfigTest(responseHeaderConfigValue string, expected interface{},
	t *testing.T) {

	// Go converts all header keys to Canonical form, hence keeping header names in canonical form here
	requestHeaderJsonMap := map[string][]string{
		"Authorization": {"hfdhfbwfbg"},
		"Content-Type":  {"application/json", "text/html", "application/pdf"},
		"X-Deviceid":    {"Device1234"},
		"X-Clientid":    {"3"},
		"X-Checksum":    {"hfsdhfbudgwq8gdqwudqu"},
	}

	actual := service.GetResponseHeaderConfigValueFromRequestHeader(responseHeaderConfigValue, requestHeaderJsonMap)

	isEqual := reflect.DeepEqual(expected, actual)
	if !isEqual {
		t.Errorf("expected value %v type %T actual value %v type %T", expected, expected,
			actual, actual)
	}
}
