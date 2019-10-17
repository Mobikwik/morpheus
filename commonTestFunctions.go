package main

import "testing"

func compareValues(expected interface{}, actual interface{}, t *testing.T) {
	switch expectedTypedValue := expected.(type) {

	case []interface{}:
		actualArrayTypeValue := actual.([]interface{})
		if len(expectedTypedValue) != len(actualArrayTypeValue) {
			t.Errorf("expected array value %v type %T actual value %v type %T", expectedTypedValue, expectedTypedValue, actualArrayTypeValue, actualArrayTypeValue)
		}
		for i, v := range expectedTypedValue {
			compareValues(v, actualArrayTypeValue[i], t)
			/*if v != actualArrayTypeValue[i] {
				t.Errorf("expected array type value %v type %T actual value %v type %T", v, v, actualArrayTypeValue[i], actualArrayTypeValue[i])
			}*/
		}

	case []string:
		actualArrayTypeValue := actual.([]string)
		if len(expectedTypedValue) != len(actualArrayTypeValue) {
			t.Errorf("expected array value %v type %T actual value %v type %T", expectedTypedValue, expectedTypedValue, actualArrayTypeValue, actualArrayTypeValue)
		}
		for i, v := range expectedTypedValue {
			compareValues(v, actualArrayTypeValue[i], t)
		}
	case map[string]interface{}:

		actualMapTypeValue := actual.(map[string]interface{})
		if len(expectedTypedValue) != len(actualMapTypeValue) {
			t.Errorf("expected array value %v type %T actual value %v type %T",
				expectedTypedValue, expectedTypedValue, actualMapTypeValue, actualMapTypeValue)
		}

		for key, v := range expectedTypedValue {
			compareValues(v, actualMapTypeValue[key], t)
			/*	if value != actualMapTypeValue[key] {
				t.Errorf("expected value %v type %T actual value %v type %T",
					value, value, actualMapTypeValue[key], actualMapTypeValue[key])
			}*/
		}

	case interface{}:
		if expected != actual {
			t.Errorf("expected value %v type %T actual value %v type %T", expected, expected, actual, actual)
		}
	default:
		t.Errorf("unexpected type passed %T", expectedTypedValue)
	}
}