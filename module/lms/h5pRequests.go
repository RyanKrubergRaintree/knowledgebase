package lms

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/utils"
)


func unzipAndUploadH5P(fileNameWithPath string) (string, error) {
	guid := strings.Replace(uuid.New().String(), "-", "", -1)
	unzipPath := utils.GetTempPath(guid + "/")

	err := utils.Unzip(fileNameWithPath, unzipPath)
	if err != nil {
		return "", err
	}
	// go through unzipped files, upload everything except dir's
	err = filepath.Walk(unzipPath,
		func(fileNameWithPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				fileNameWithoutTempPath := strings.Replace(fileNameWithPath, utils.GetTempPath(""), "", -1)
				s3Path := filepath.FromSlash("H5P/lessons/" + fileNameWithoutTempPath)
				s3Path = strings.Replace(s3Path, string(filepath.Separator), "/", -1) // fix path for S3
				_, err:= uploadSingleFileToS3(s3Path, fileNameWithPath, "")
				if err != nil {
					return err
				}
			}
			return nil
		})
	_ = os.RemoveAll(unzipPath)
	if err != nil {
		return "", err
	}

	// upload template.html as it's needed to show the H5P content
	workingDir, _ := os.Getwd()
	fileNameWithPath = filepath.FromSlash(workingDir + "/client/H5Ptemplate.html")
	return uploadSingleFileToS3("H5P/lessons/"+guid+"/template.html", fileNameWithPath, "")
}


func listLessonsFromBucket(w http.ResponseWriter) {
	s3Client, err := getS3Client()
	if err != nil {
		kb.WriteResult(w, err)
		return
	}

	bucket := utils.GetEnvWithDefault("AWS_KB_BUCKET", kb.DefaultBucketName)
	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String("H5P/lessons"),
	}

	var result struct {
		Lessons []string `json:"lessons"`
	}

	err = s3Client.ListObjectsPages(params,
		func(response *s3.ListObjectsOutput, lastPage bool) bool {
			// Match URL-s up to lesson ID
			re := regexp.MustCompile(`^.+([/]{2}).+?([/]{1}).+?([/]{1}).+?([/]{1}).+?([/]{1})`)
			lessonLink := ""

			for _, item := range response.Contents {
				temp := re.FindString("https://" + bucket + ".s3.amazonaws.com/" + *item.Key)
				if lessonLink != temp {
					lessonLink = temp
					result.Lessons = append(result.Lessons, lessonLink+"template.html")
				}
			}
			// continue with the next page
			return true
		})

	if err != nil {
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



