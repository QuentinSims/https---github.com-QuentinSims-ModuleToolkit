package toolkit

import "testing"

func TestTools_RandomStrings(t *testing.T) {
	// Mock the Tools struct with desired configuration
	mockTools := Tools{
		MaxFileSize:      0,          // Set max file size to 0 for testing purposes
		AllowedFileTypes: []string{}, // Set allowed file types to empty for testing purposes
	}

	// Specify options for random string generation
	options := RandomStringOptions{
		IncludeLetters: true,
		IncludeNumbers: true,
		Length:         10,
	}

	// Generate random string using the mock Tools struct and options
	s := mockTools.RandomString(options)

	// Check if the generated string has correct length
	if len(s) != options.Length {
		t.Errorf("Wrong length of random string: expected %d, got %d", options.Length, len(s))
	}
}
