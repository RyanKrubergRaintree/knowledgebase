package lms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/utils"
)

const timeout = 60 * 60 * time.Second // max time for single upload (1h)

// Uploads single video file from the server; Returns S3 path if successful
func uploadVideoFileFromServerToS3(fileNameWithPath, clientID, environment, guid string) (string, error) {
	year := strconv.Itoa(time.Now().Year())
	path := "videos/" + environment + "/" + clientID + "/" + year + "/" + guid + "_" + filepath.Base(fileNameWithPath)

	return uploadSingleFileToS3(path, fileNameWithPath, "rt-kb-videos")
}

// Deletes single file from S3
func deleteFileFromS3(key, bucket string) error {
	session, err := session.NewSession(&aws.Config{Region: aws.String(utils.GetEnvWithDefault("AWS_REGION", kb.DefaultRegion))})
	if err != nil {
		return kb.ErrUnableToCreateS3Session
	}
	s3Client := s3.New(session)

	if bucket == "" {
		bucket = utils.GetEnvWithDefault("AWS_KB_BUCKET", kb.DefaultBucketName)
	}

	prefix := "https://" + bucket + ".s3.amazonaws.com/"
	key = strings.Replace(key, prefix, "", -1)

	if ok := doesFileExistsInS3(s3Client, bucket, key); !ok {
		return kb.ErrDoesNotExist
	}

	_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)})
	if err != nil {
		return kb.ErrUnableToDelete
	}


	err = s3Client.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return kb.ErrUnableToDelete
	}

	return nil
}

// Uploads single file from the server; Returns S3 path if successful
func uploadFileFromServerToS3(fileNameWithPath string) (string, error) {
	fileExtension := strings.ToUpper(filepath.Ext(fileNameWithPath))

	if fileExtension == ".H5P" {
		return unzipAndUploadH5P(fileNameWithPath)
	}
	return uploadSingleFileToS3("", fileNameWithPath, "")
}


// Uploads single file from the server; Returns S3 path if successful
// S3 full path can be specified (optional)
func uploadSingleFileToS3(destinations3Path, fileNameWithPath, bucket string) (string, error) {
	var key *string
	var uploadedFilePath string
	fileName := filepath.Base(fileNameWithPath)
	if bucket == "" {
		bucket = utils.GetEnvWithDefault("AWS_KB_BUCKET", kb.DefaultBucketName)
	}

	if destinations3Path == "" {
		key = aws.String(fileName) // upload to Root dir if no specific path given
	} else {
		key = aws.String(destinations3Path)
	}
	uploadedFilePath = "https://" + bucket + ".s3.amazonaws.com/" + *key

	defaultRegion := utils.GetEnvWithDefault("AWS_REGION", kb.DefaultRegion)
	// Init session and service. Uses ENV variables AWS_ACCESS_KEY_ID & AWS_SECRET_ACCESS_KEY
	session, err := session.NewSession(&aws.Config{Region: aws.String(defaultRegion)})
	if err != nil {
		return "", err
	}
	svc := s3.New(session)

	// To abort the upload if it takes more than timeout seconds
	ctx, cancelFn := context.WithTimeout(context.Background(), timeout)
	defer cancelFn() // Ensure the context is canceled to prevent leaking.

	file, IOerr := os.Open(fileNameWithPath)
	if IOerr != nil {
		return "", IOerr
	}
	defer file.Close()

	_, err = svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         key,
		Body:        file,
		ContentType: getContentType(fileNameWithPath),
	})

	if err != nil {
		return "", err
	}

	return uploadedFilePath, nil
}

// Grape JS image uploading to [bucket]/customer/database/year/[filename];
// Returns S3 path if request was successful
func uploadImageFromServerToS3(clientname, database, filename, sourceFile string) (string, error) {
	if strings.TrimSpace(clientname) == "" || strings.TrimSpace(database) == "" || strings.TrimSpace(sourceFile) == "" || strings.TrimSpace(filename) == "" {
		return "", errors.New("client name, database, file name or file data missing")
	}

	year := strconv.Itoa(time.Now().Year())
	path := "customers/" + clientname + "/" + database + "/" + year + "/" + filename

	return uploadSingleFileToS3(path, sourceFile, "")
}


// bucket access needs to be verified; although listing files != file access
func listFilesForGivenBucketAndPath(bucket, path string, w http.ResponseWriter) {
	if bucket == "" {
		bucket = utils.GetEnvWithDefault("AWS_KB_BUCKET", kb.DefaultBucketName)
	}

	defaultRegion := utils.GetEnvWithDefault("AWS_REGION", kb.DefaultRegion)

	// Init session and service. Uses ENV variables AWS_ACCESS_KEY_ID & AWS_SECRET_ACCESS_KEY
	session, err1 := session.NewSession(&aws.Config{Region: aws.String(defaultRegion)})
	if err1 != nil {
		fmt.Fprintf(w, "Unable to list items from bucket %q, %v", bucket, err1)
		kb.WriteResult(w, err1)
		return
	}
	svc := s3.New(session)

	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(path),
	}

	var result struct {
		Uploads []string `json:"uploads"`
	}

	err := svc.ListObjectsPages(params,
		func(response *s3.ListObjectsOutput, lastPage bool) bool {

			for _, item := range response.Contents {
				result.Uploads = append(result.Uploads, "https://"+bucket+".s3.amazonaws.com/"+*item.Key)
			}
			// continue with the next page
			return true
		})

	if err != nil {
		fmt.Fprintf(w, "Unable to list all items from bucket %q, %v", bucket, err)
		kb.WriteResult(w, err)
		return
	}

	data, err := json.Marshal(result)
	if err != nil {
		kb.WriteResult(w, err)
		return
	}

	w.Write(data)
	w.Header().Set("Content-Type", "application/json")
}




