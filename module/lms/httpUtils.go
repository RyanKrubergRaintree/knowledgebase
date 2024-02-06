package lms

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/utils"
)


func getContentType(fileNameWithPath string) *string {

	fileExtension := strings.ToUpper(filepath.Ext(fileNameWithPath))
	if fileExtension == ".HTML" {
		return aws.String("text/html")
	} else if fileExtension == ".CSS" {
		return aws.String("text/css")
	} else if fileExtension == ".JS" {
		return aws.String("text/javascript")
	} else if fileExtension == ".JSON" {
		return aws.String("application/json")
	} else {
		file, err := os.Open(fileNameWithPath)
		if err != nil {
			return aws.String("binary/octet-stream")
		}
		defer file.Close()

		header := make([]byte, 512)
		_, _ = file.Read(header)
		return aws.String(http.DetectContentType(header))
	}
}

// Saves single(first) file from http request to temp folder. Expects form key to be "file".
// on success returns full path to the received file
func saveFileFromHttpRequestToServer(r *http.Request) (string, error) {
	file, fileHeader, fileErr := r.FormFile("file")

	if fileErr != nil {
		log.Println(fileErr)
		log.Println("Content type was: " + r.Header.Get("Content-type"))
		return "", fileErr
	}
	if file == nil || fileHeader == nil {
		return "", kb.ErrBadRequest
	}

	fileNameWithPath := utils.GetTempPath(fileHeader.Filename)

	f, err := os.Create(fileNameWithPath)
	if err != nil {
		return "", err
	}
	//nolint:errcheck
	io.Copy(f, file)
	defer f.Close()

	log.Println("HTTP -> Server upload complete. Received file: " + fileHeader.Filename)
	return fileNameWithPath, nil
}

// Grape JS; TODO: simultaneous uploads with same file name?
func saveImageFileFromHttpRequestToServer(r *http.Request) (string, error) {
	image := r.PostFormValue("image")
	filename := r.PostFormValue("filename")

	if image == "" || filename == "" {
		return "", kb.ErrBadRequest
	}

	decoded, err := base64.StdEncoding.DecodeString(image)
	if err != nil {
		return "", err
	}

	fileNameWithPath := utils.GetTempPath(filename)

	f, err := os.Create(fileNameWithPath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	//nolint:errcheck
	f.Write(decoded)

	return fileNameWithPath, nil
}

