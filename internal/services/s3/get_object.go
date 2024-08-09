package s3_services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rodrigoenzohernandez/transactions-processor/internal/utils/logger"
)

var log = logger.GetLogger("s3_services")

// Gets the object from the specified s3 bucket and key and returns it
func GetObject(bucket string, key string, s3Client *s3.S3) (*s3.GetObjectOutput, error) {

	object, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Error(fmt.Sprintf("Unable to download object %q, %v", key, err))
		return nil, err
	}

	return object, nil
}
