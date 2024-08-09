package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/rodrigoenzohernandez/transactions-processor/internal/utils"
	"github.com/rodrigoenzohernandez/transactions-processor/internal/utils/logger"
)

type Event struct {
	Records []Record `json:"records"`
}

type Record struct {
	S3 S3 `json:"s3"`
}

type S3 struct {
	Bucket Bucket `json:"bucket"`
	Object Object `json:"object"`
}

type Bucket struct {
	Name string
}

type Object struct {
	Key string
}

var log = logger.GetLogger("files_processor_lambda")

func handler(ctx context.Context, event events.S3Event) {
	eventJson, _ := json.Marshal(event)
	var data Event
	json.Unmarshal(eventJson, &data)
	bucket := data.Records[0].S3.Bucket.Name
	key := data.Records[0].S3.Object.Key

	log.Info(fmt.Sprintf("An object was uploaded to bucket %s with key %s", bucket, key))

	session := session.Must(session.NewSession())

	s3Client := s3.New(session)

	// create service to get object

	object, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Error(fmt.Sprintf("Unable to download object %q, %v", key, err))
		return
	}
	defer object.Body.Close()

	records := utils.GetRecordsFromBuffer(object.Body)

	report, _ := json.Marshal(utils.GenerateReport(records))

	sqsClient := sqs.New(session)

	// create service GetQueueUrl
	resultURL, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String("reports_queue"),
	})
	if err != nil {
		log.Error(fmt.Sprintf("Failed to get queue URL: %v", err))
		return
	}

	queueURL := *resultURL.QueueUrl

	log.Info(fmt.Sprintf("Records: %s", records))

	// create service SendMessage

	resultSend, err := sqsClient.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(report)),
		QueueUrl:    aws.String(queueURL),
	})
	if err != nil {
		log.Error(fmt.Sprintf("Failed to send message to SQS queue %v", err))

		return
	}

	log.Info(fmt.Sprintf("Successfully sent message to SQS queue with ID: %s\n", *resultSend.MessageId))

}

func main() {
	lambda.Start(handler)
}
