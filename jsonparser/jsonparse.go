package jsonparser

import "errors"

type Json struct {
}

func handleParseObject(j string) (string, error) {
	if j == "" {
		myError := errors.New("not valid json, jusr an empty string")
		return "", myError
	}

	var abcd string
	return string(abcd), nil
}
