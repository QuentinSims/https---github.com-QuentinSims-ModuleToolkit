package jsonparser

import "errors"

type Json struct {
}

func handleParseObject(jsonData []byte) (bool, error) {
	if len(jsonData) == 0 {
		myError := errors.New("valid json, but empty")
		return false, myError
	}

	return true, nil
}
