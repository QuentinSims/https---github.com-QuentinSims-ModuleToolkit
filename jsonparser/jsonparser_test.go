package jsonparser

import "testing"

func TestJsonValidJson_handleParseObject(t *testing.T) {
	jsonData := []byte(`{"key": "value"}`)
	//expectedResult := true

	result, err := handleParseObject(jsonData)

	if !result {
		t.Errorf("Expected true: %t", result)
	}

	if err != nil {
		t.Errorf("Expected no error, but got: %s", err.Error())
	}
}

func TestJsonInvalidJson_handleParseObject(t *testing.T) {
	invalidJson := []byte(`{"key": value}`)
	expectedErrorMessage := "invalid character 'v' looking for beginning of value"

	result, err := handleParseObject(invalidJson)

	if result {
		t.Errorf("Expected false: %t", result)
	}

	if err == nil {
		t.Error("Expected an error, but got nil")
	} else if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message: %s, got: %s", expectedErrorMessage, err.Error())
	}
}
