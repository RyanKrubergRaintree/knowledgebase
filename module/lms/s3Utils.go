package lms

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/utils"
)

func getSignedLink(key, bucket string,  w http.ResponseWriter) (string, error) {
	s3Client, err := getS3Client()
	if err != nil {
		return "", err
	}

	if bucket == "" {
		bucket = utils.GetEnvWithDefault("AWS_KB_BUCKET", kb.DefaultBucketName)
	}

	prefix := "https://" + bucket + ".s3.amazonaws.com/"
	key = strings.Replace(key, prefix, "", -1)

	req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	urlStr, err := req.Presign(8 * 60 * time.Minute)

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString([]byte(urlStr)), nil
}

func doesFileExistsInS3(s3Client *s3.S3, bucket string, key string) bool {
    _, err := s3Client.HeadObject(&s3.HeadObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(key),
    })

	if err != nil {
		return false
	}

	return true
}

// todo: error out only if bucket does not exist and err. happens; i.e ignore bucket exists errors
func createBucket(bucketName string) error {
	s3Client, err := getS3Client()
	if err != nil {
		return  err
	}

	_, err = s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return err
	}

	err = s3Client.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return err
	}

	return nil
}


// creates a new S3 client and returns it along with any error.
func getS3Client() (*s3.S3, error) {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(utils.GetEnvWithDefault("AWS_REGION", kb.DefaultRegion)),
	})
	if err != nil {
		return nil, kb.ErrUnableToCreateS3Session
	}

	return s3.New(session), nil
}
