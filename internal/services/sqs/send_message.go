package sqs_services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// Sends a SQS message to a queue
func SendMessage(message []byte, queueUrl string, sqsClient *sqs.SQS) error {

	resultSend, err := sqsClient.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(message)),
		QueueUrl:    aws.String(queueUrl),
	})
	if err != nil {
		log.Error(fmt.Sprintf("Failed to send message to SQS queue %v", err))
		return err
	}

	log.Info(fmt.Sprintf("Successfully sent message to SQS queue with ID: %s\n", *resultSend.MessageId))

	return nil

}
