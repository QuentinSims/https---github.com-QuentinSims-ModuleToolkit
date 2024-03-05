package jsonparser

import "testing"

func TestJsonEmptyorNil_handleParseObject(t *testing.T) {
	var emptyString []byte
	expectedErrorMessage := "valid json, but empty"

	result, err := handleParseObject(emptyString)

	if result {
		t.Errorf("Expected true: %s", result)
	}

	if err == nil {
		t.Error("Expected an error, but got nil")
	} else if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message: %s, got: %s", expectedErrorMessage, err.Error())
	}
}
