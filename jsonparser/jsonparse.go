package jsonparser

import (
	"encoding/json"
	"errors"
)

type Json struct {
	Data interface{}
}

func handleParseObject(jsonData []byte) (bool, error) {
	if len(jsonData) == 0 {
		myError := errors.New("valid JSON, but empty")
		return false, myError
	}

	var parsedData interface{}
	err := json.Unmarshal(jsonData, &parsedData)
	if err != nil {
		return false, err
	}

	return true, nil
}
