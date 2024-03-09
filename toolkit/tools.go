package toolkit

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//const randomStringSource = "abcdefghijklmnopqrstuwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_+"

// Tools is the type used to instantiate this module.
// any variable of this type has access to all the methods with the
// reciever *Tools
type Tools struct {
	MaxFileSize       int
	AllowedFileTypes  []string
	MaxJSONSize       int
	AllowUnkownFields bool
}

// Define character sets
const (
	letters      = "abcdefghijklmnopqrstuvwxyz"
	capitals     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers      = "0123456789"
	specialChars = "_+!@#$%^&*?/"
	allChars     = letters + capitals + numbers + specialChars
)

// RandomStringOptions defines options for customizing random string generation
type RandomStringOptions struct {
	IncludeLetters  bool
	IncludeCapitals bool
	IncludeNumbers  bool
	IncludeSpecial  bool
	Length          int
}

// random string returns a string of random characters of length of n, using randomStringSource
// as the source for the string
// RandomString returns a string of random characters based on the specified options
func (t *Tools) RandomString(options RandomStringOptions) string {
	var sourceChars string

	if options.IncludeLetters {
		sourceChars += letters
	}
	if options.IncludeCapitals {
		sourceChars += capitals
	}
	if options.IncludeNumbers {
		sourceChars += numbers
	}
	if options.IncludeSpecial {
		sourceChars += specialChars
	}

	// If no specific characters are included, default to allChars
	if sourceChars == "" {
		sourceChars = allChars
	}

	// Generate random string
	var result strings.Builder
	for i := 0; i < options.Length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(sourceChars))))
		if err != nil {
			// handle error
			errorMessage := fmt.Sprintf("Error generating random string: %v", err)
			return errorMessage
		}
		result.WriteByte(sourceChars[n.Int64()])
	}

	return result.String()
}

// uploaded file is a struct used to save information about an uploaded file
type UploadedFile struct {
	NewFileName      string
	OriginalFileName string
	FileSize         int64
}

// UploadOneFile is a method of the Tools struct responsible for handling one uploaded file.
// It takes an *http.Request containing the uploaded file, the directory to which the file will be saved,
// and an optional boolean flag indicating whether to rename the file.
// It returns a slice of UploadedFile pointers representing the uploaded file and an error if any occurred during the handling process.
func (t *Tools) UploadOneFile(r *http.Request, uploadDir string, rename ...bool) (*UploadedFile, error) {
	renameFile := true

	if len(rename) > 0 {
		renameFile = rename[0]
	}

	files, err := t.UploadFiles(r, uploadDir, renameFile)
	if err != nil {
		return nil, err
	}

	return files[0], nil
}

// handleUploadedFile is a helper function responsible for handling an individual uploaded file.
// It takes a multipart.FileHeader representing the uploaded file, the directory to which the file will be saved,
// and a boolean flag indicating whether to rename the file.
// It returns an UploadedFile pointer representing the uploaded file and an error if any occurred during the handling process.
func handleUploadedFile(t *Tools, hdr *multipart.FileHeader, uploadDir string, renameFile bool) (*UploadedFile, error) {
	var uploadedFile UploadedFile

	// Open the uploaded file
	infile, err := hdr.Open()
	if err != nil {
		return nil, err
	}
	defer infile.Close()

	// Read a small buffer from the beginning of the file to detect its content type
	buff := make([]byte, 512)
	_, err = infile.Read(buff)
	if err != nil {
		return nil, err
	}

	// Check if the file type is permitted
	allowed := false
	fileType := http.DetectContentType(buff)

	// If allowed file types are specified, check if the detected file type is in the allowed list
	if len(t.AllowedFileTypes) > 0 {
		for _, x := range t.AllowedFileTypes {
			if strings.EqualFold(fileType, x) {
				allowed = true
			}
		}
	} else {
		// If no allowed file types are specified, consider all file types as allowed
		allowed = true
	}

	// If the file type is not permitted, return an error
	if !allowed {
		return nil, errors.New("the uploaded file type is not permitted")
	}

	// Reset the file read pointer to the beginning
	_, err = infile.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	// Determine the new file name based on the renaming flag
	if renameFile {
		options := RandomStringOptions{
			IncludeLetters: true,
			IncludeNumbers: true,
			Length:         25,
		}
		uploadedFile.NewFileName = fmt.Sprintf("%s%s", t.RandomString(options), filepath.Ext(hdr.Filename))
	} else {
		uploadedFile.NewFileName = hdr.Filename
	}

	uploadedFile.OriginalFileName = hdr.Filename

	// Create the output file in the specified upload directory
	var outfile *os.File
	defer outfile.Close()
	if outfile, err = os.Create(filepath.Join(uploadDir, uploadedFile.NewFileName)); err != nil {
		return nil, err
	} else {
		// Copy the file contents to the output file and record the file size
		fileSize, err := io.Copy(outfile, infile)
		if err != nil {
			return nil, err
		}
		uploadedFile.FileSize = fileSize
	}

	// Return the uploaded file information
	return &uploadedFile, nil
}

// UploadFiles is a method of the Tools struct responsible for handling multiple uploaded files.
// It takes an *http.Request containing the uploaded files, the directory to which the files will be saved,
// and an optional boolean flag indicating whether to rename the files.
// It returns a slice of UploadedFile pointers representing the uploaded files and an error if any occurred during the handling process.
func (t *Tools) UploadFiles(r *http.Request, uploadDir string, rename ...bool) ([]*UploadedFile, error) {
	renameFile := true
	if len(rename) > 0 {
		renameFile = rename[0]
	}

	// Initialize a slice to store information about uploaded files
	var uploadedFiles []*UploadedFile

	// Set a default maximum file size if not specified
	if t.MaxFileSize == 0 {
		t.MaxFileSize = 1024 * 1024 * 1024
	}

	// Parse the multipart form data from the HTTP request
	err := r.ParseMultipartForm(int64(t.MaxFileSize))
	if err != nil {
		return nil, errors.New("the uploaded file is too big")
	}

	// Iterate through each file header in the multipart form data
	for _, fHeaders := range r.MultipartForm.File {
		for _, hdr := range fHeaders {
			// Handle the uploaded file using the helper function
			uploadedFile, err := handleUploadedFile(t, hdr, uploadDir, renameFile)
			if err != nil {
				return uploadedFiles, err
			}
			// Append information about the uploaded file to the slice
			uploadedFiles = append(uploadedFiles, uploadedFile)
		}
	}

	// Return the slice of uploaded files
	return uploadedFiles, nil
}

type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (t *Tools) ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxByte := 1024 * 1024

	if t.MaxJSONSize != 0 {
		maxByte = t.MaxJSONSize
	}

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxByte))

	dec := json.NewDecoder(r.Body)

	if !t.AllowUnkownFields {
		dec.DisallowUnknownFields()
	}

	err := dec.Decode(data)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshallTypeError *json.UnmarshalTypeError
		var invalidUnmarshallError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed Json")

		case errors.As(err, &unmarshallTypeError):
			if unmarshallTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect Json type for field %q", unmarshallTypeError)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshallTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json:unknow field")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d", maxByte)

		case errors.As(err, &invalidUnmarshallError):
			return fmt.Errorf("error unmarshalling JSON: %s", err.Error())

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must contain only one Json value")
	}

	return nil
}
