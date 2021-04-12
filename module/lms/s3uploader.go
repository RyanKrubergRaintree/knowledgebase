package lms

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

const timeout = 60 * 60 * time.Second // max time for single upload (1h)

// Uploads single file from the server; Returns S3 path if successful
func uploadFileFromServerToS3(fileNameWithPath string) (error, string) {
	fileExtension := strings.ToUpper(filepath.Ext(fileNameWithPath))

	if fileExtension == ".H5P" {
		return unzipAndUploadH5P(fileNameWithPath)
	}
	return uploadSingleFileToS3("", fileNameWithPath)
}

func unzipAndUploadH5P(fileNameWithPath string) (error, string) {
	guid := strings.Replace(uuid.New().String(), "-", "", -1)
	unzipPath := getTempPath(guid + "/")

	err := Unzip(fileNameWithPath, unzipPath)
	if err != nil {
		return err, ""
	}
	// go through unzipped files, upload everything except dir's
	err = filepath.Walk(unzipPath,
		func(fileNameWithPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				fileNameWithoutTempPath := strings.Replace(fileNameWithPath, getTempPath(""), "", -1)
				s3Path := filepath.FromSlash("H5P/lessons/" + fileNameWithoutTempPath)
				s3Path = strings.Replace(s3Path, string(filepath.Separator), "/", -1) // fix path for S3
				err, _ := uploadSingleFileToS3(s3Path, fileNameWithPath)
				if err != nil {
					return err
				}
			}
			return nil
		})
	_ = os.RemoveAll(unzipPath)
	if err != nil {
		return err, ""
	}

	// upload template.html as it's needed to show the H5P content
	workingDir, _ := os.Getwd()
	fileNameWithPath = filepath.FromSlash(workingDir + "/client/H5Ptemplate.html")
	return uploadSingleFileToS3("H5P/lessons/"+guid+"/template.html", fileNameWithPath)
}

// Uploads single file from the server; Returns S3 path if successful
// S3 full path can be specified (optional)
func uploadSingleFileToS3(destinations3Path string, fileNameWithPath string) (error, string) {
	var key *string
	var uploadedFilePath string
	fileName := filepath.Base(fileNameWithPath)
	bucket := getEnvWithDefault("AWS_KB_BUCKET", "rt-knowledge-base-dev")

	if destinations3Path == "" {
		key = aws.String(fileName) // upload to Root dir if no specific path given
	} else {
		key = aws.String(destinations3Path)
	}
	uploadedFilePath = "https://" + bucket + ".s3.amazonaws.com/" + *key

	defaultRegion := getEnvWithDefault("AWS_REGION", "us-east-1")
	// Init session and service. Uses ENV variables AWS_ACCESS_KEY_ID & AWS_SECRET_ACCESS_KEY
	sess, err := session.NewSession(&aws.Config{Region: aws.String(defaultRegion)})
	if err != nil {
		return err, ""
	}
	amazonFileBucket := s3.New(sess)

	// To abort the upload if it takes more than timeout seconds
	ctx, cancelFn := context.WithTimeout(context.Background(), timeout)
	defer cancelFn() // Ensure the context is canceled to prevent leaking.

	file, IOerr := os.Open(fileNameWithPath)
	if IOerr != nil {
		return IOerr, ""
	}
	defer file.Close()

	_, err = amazonFileBucket.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         key,
		Body:        file,
		ContentType: getContentType(fileNameWithPath),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			return err, "" // timeout
		} else {
			return err, ""
		}
	}

	return nil, uploadedFilePath
}

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

func getEnvWithDefault(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func getTempPath(append string) string {
	workingDir, _ := os.Getwd()
	workingDir += "/temp/" + append

	return filepath.FromSlash(workingDir)
}

// todo: return JSON; Filter for H5PDist folder & list only uploads i.e guid/template.hmtl files; 500 response on errors
func ListFilesFromBucket(w http.ResponseWriter) {
	bucket := getEnvWithDefault("AWS_KB_BUCKET", "rt-knowledge-base-dev")
	defaultRegion := getEnvWithDefault("AWS_REGION", "us-east-1")

	// Init session and service. Uses ENV variables AWS_ACCESS_KEY_ID & AWS_SECRET_ACCESS_KEY
	sess, err1 := session.NewSession(&aws.Config{Region: aws.String(defaultRegion)})
	if err1 != nil {
		fmt.Fprintf(w, "Unable to list items in bucket %q, %v", bucket, err1)
		return
	}
	amazonFileBucket := s3.New(sess)

	resp, err := amazonFileBucket.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	if err != nil {
		fmt.Fprintf(w, "Unable to list items in bucket %q, %v", bucket, err)
		return
	}

	for _, item := range resp.Contents {
		fmt.Fprintf(w, "https://"+bucket+".s3.amazonaws.com/"+*item.Key)
	}
}

// Saves single(first) file from http request to temp folder. Expects form key to be "file".
// on success returns full path to the received file
// todo: return 400 in case input was invalid
func saveFileFromHttpRequestToServer(r *http.Request) (error, string) {
	file, fileHeader, fileErr := r.FormFile("file")

	if fileErr != nil {
		log.Println(fileErr)
		log.Println("Content type was: " + r.Header.Get("Content-type"))
		return fileErr, ""
	}
	if file == nil || fileHeader == nil {
		return errors.New("Upload error: file and/or header missing."), ""
	}

	fileNameWithPath := getTempPath(fileHeader.Filename)

	f, err := os.Create(fileNameWithPath)
	if err != nil {
		return err, ""
	}
	io.Copy(f, file)
	defer f.Close()

	log.Println("HTTP -> Server upload complete. Received file: " + fileHeader.Filename)
	return nil, fileNameWithPath
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	os.MkdirAll(dest, 0644)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)
		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
