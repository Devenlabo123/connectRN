package api

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestJpegToPNGHappyPath(t *testing.T) {
	h, err := CreateImageHandler()
	if err != nil {
		t.Fatalf("error creating handler")
	}

	// read test file into request body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("file", "JPEG_example_flower.png")
	if err != nil {
		t.Fatalf("error handling test file")
	}

	file, err := os.Open("upload-test.png")
	if err != nil {
		t.Fatalf("error handling test file")
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		t.Fatalf("error handling test file")
	}

	writer.Close()


	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/images", bytes.NewReader(body.Bytes()))
	request.Header.Set("Content-Type", writer.FormDataContentType())

	h.ServeHTTP(responseRecorder, request)

	assert.Equal(t, "image/png", responseRecorder.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	responseBodyBytes, err := io.ReadAll(responseRecorder.Body)

	// verify image type
	actualContentType := http.DetectContentType(responseBodyBytes)
	assert.Equal(t, "image/png", actualContentType)


	// verify image size
	image, _, err := image.Decode(bytes.NewReader(responseBodyBytes))
	assert.Equal(t, "(256,256)", image.Bounds().Max.String())
}

func TestPNGToPNGHappyPath(t *testing.T) {
	h, err := CreateImageHandler()
	if err != nil {
		t.Fatalf("error creating handler")
	}

	// read test file into request body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("file", "upload-test.png")
	if err != nil {
		t.Fatalf("error handling test file")
	}

	file, err := os.Open("upload-test.png")
	if err != nil {
		t.Fatalf("error handling test file")
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		t.Fatalf("error handling test file")
	}

	writer.Close()


	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/images", bytes.NewReader(body.Bytes()))
	request.Header.Set("Content-Type", writer.FormDataContentType())

	h.ServeHTTP(responseRecorder, request)

	assert.Equal(t, "image/png", responseRecorder.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	responseBodyBytes, err := io.ReadAll(responseRecorder.Body)

	// verify image type
	actualContentType := http.DetectContentType(responseBodyBytes)
	assert.Equal(t, "image/png", actualContentType)


	// verify image size
	image, _, err := image.Decode(bytes.NewReader(responseBodyBytes))
	assert.Equal(t, "(256,256)", image.Bounds().Max.String())
}