package toolkit

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
)

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

var uploadTests = []struct {
	name          string
	allowedTypes  []string
	renameFile    bool
	errorExpected bool
}{
	{name: "allowed no rename", allowedTypes: []string{"image/jpeg", "image/png"}, renameFile: false, errorExpected: false},
}

func TestTools_UploadFiles(t *testing.T) {
	for _, e := range uploadTests {
		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)

		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer writer.Close()
			defer wg.Done()

			part, err := writer.CreateFormFile("file", "./testdata/img.png")
			if err != nil {
				t.Error(err)
			}

			f, err := os.Open("./testdata/img.png")
			if err != nil {
				t.Error(err)
			}
			defer f.Close()

			img, _, err := image.Decode(f)
			if err != nil {
				t.Error("error decoding image", err)
			}

			err = png.Encode(part, img)
			if err != nil {
				t.Error(err)
			}
		}()

		request := httptest.NewRequest("POST", "/", pr)
		request.Header.Add("Content-Type", writer.FormDataContentType())

		var testTools Tools
		testTools.AllowedFileTypes = e.allowedTypes

		uploadedFiles, err := testTools.UploadFiles(request, "./testdata/uploads", e.renameFile)
		if err != nil && !e.errorExpected {
			t.Error(err)
		}

		if !e.errorExpected {
			if _, err := os.Stat(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].NewFileName)); os.IsNotExist(err) {
				t.Errorf("%s: expected file to exist: %s", e.name, err.Error())
			}

			_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].NewFileName))
		}

		if !e.errorExpected && err != nil {
			t.Errorf("%s:error expected but none recieved", e.name)
		}

		wg.Wait()

	}
}

var jsonTests = []struct {
	name          string
	json          string
	errorExpected bool
	maxSize       int
	allowUnknown  bool
}{
	{name: "good json", json: `{"Quentin": "Sims"}`, errorExpected: false, maxSize: 1024, allowUnknown: false},
}

func TestTools_ReadJSON(t *testing.T) {

	var testTool Tools

	for _, e := range jsonTests {
		testTool.MaxJSONSize = e.maxSize

		testTool.AllowUnkownFields = e.allowUnknown

		var decodedJson struct {
			Quentin string `json:"Quentin"`
		}

		req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(e.json)))
		if err != nil {
			t.Log("Error: ", err)
		}

		rr := httptest.NewRecorder()

		err = testTool.ReadJSON(rr, req, &decodedJson)

		if e.errorExpected && err == nil {
			t.Errorf("%s:error expected, but none received", e.name)
		}

		if !e.errorExpected && err != nil {
			t.Errorf("%s:error not expected, but none received: %s", e.name, err.Error())
		}

		req.Body.Close()
	}
}
