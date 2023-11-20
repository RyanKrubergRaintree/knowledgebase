package kb

import "errors"

// AWS S3 related
const DefaultBucketName = "rt-knowledge-base-dev"
const DefaultRegion = "us-east-1"


var (
	// S3 requests related errors
	ErrUnableToDelete          = errors.New("Unable to delete given object.")
	ErrUnableToCreateS3Session = errors.New("Unable to create S3 session.")
	ErrDoesNotExist            = errors.New("404 Not Found")
	ErrBadRequest              = errors.New("Bad request, likely due to invalid input.")
)



