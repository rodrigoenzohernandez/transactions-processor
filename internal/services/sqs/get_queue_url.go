package sqs_services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/rodrigoenzohernandez/transactions-processor/internal/utils/logger"
)

var log = logger.GetLogger("sqs_services")

// Gets the queueURL searching by queueName
func GetQueueURL(queueName string, sqsClient *sqs.SQS) (string, error) {
	resultURL, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		log.Error(fmt.Sprintf("Failed to get queue URL: %v", err))
		return "", err
	}

	queueURL := *resultURL.QueueUrl

	return queueURL, nil
}
