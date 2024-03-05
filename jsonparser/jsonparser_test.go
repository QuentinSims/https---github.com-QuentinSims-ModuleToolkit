package jsonparser

import "testing"

func TestJsonEmptyorNil_handleParseObject(t *testing.T) {
	emptyString := ""
	expectedErrorMessage := "not valid json, jusr an empty string"

	result, err := handleParseObject(emptyString)

	if result != "" {
		t.Errorf("Expected empty string as result, got: %s", result)
	}

	if err == nil {
		t.Error("Expected an error, but got nil")
	} else if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message: %s, got: %s", expectedErrorMessage, err.Error())
	}
}
