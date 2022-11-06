package api

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"

	"github.com/nfnt/resize"
)

type ImageHandler struct {}

func CreateImageHandler() (*ImageHandler, error) {
	return &ImageHandler{}, nil
}

func (i ImageHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		// validate auth

		err := i.handleCreateUserRequest(writer, request)
		if err != nil {
			log.Error(err)
			WriteStatusCode(writer, http.StatusBadRequest)
			WriteResponse(writer, err.Error())
		}
	default:
		WriteStatusCode(writer, http.StatusMethodNotAllowed)
	}
}

func (i ImageHandler) handleCreateUserRequest (writer http.ResponseWriter, request *http.Request) error {
	// set max size of upload
	request.ParseMultipartForm(10 << 20)

	// read uploaded file into memory as byte[]
	file, handler, err := request.FormFile("file")
	if err != nil {
		log.Error("Error Retrieving the File")
		return fmt.Errorf(err.Error())
	}

	defer file.Close()
	log.Debug("Uploaded File: %+v\n", handler.Filename)
	log.Debug("File Size: %+v\n", handler.Size)
	log.Debug("MIME Header: %+v\n", handler.Header)

	uploadedFileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	// convert file to png
	pngImageBytes, err := convertFileToPng(uploadedFileBytes)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	// resize image
	resizedPngImageBytes, err := resizeImage(pngImageBytes)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	writer.Header().Set("Content-Type", "image/png")
	writer.Write(resizedPngImageBytes)

	return nil
}

func convertFileToPng(imageBytes []byte) ([]byte, error) {
	contentType := http.DetectContentType(imageBytes)

	switch contentType {
	case "image/png":
		log.Info("file already is a png")
		return imageBytes, nil
	case "image/jpeg":
		img, err := jpeg.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return nil, fmt.Errorf("unable to decode jpeg: " + err.Error())
		}

		buf := new(bytes.Buffer)
		if err := png.Encode(buf, img); err != nil {
			return nil, fmt.Errorf("unable to encode png: " + err.Error())
		}

		return buf.Bytes(), nil
	}

	return nil, fmt.Errorf("unable to convert %#v to png", contentType)
}

func resizeImage(pngImageBytes []byte) ([]byte, error) {
	pngImage, _, err := image.Decode(bytes.NewReader(pngImageBytes))
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	resized := resize.Resize(256, 256, pngImage, resize.Lanczos3)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, resized, nil)

	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return buf.Bytes(), nil
}
