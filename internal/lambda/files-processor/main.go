package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/rodrigoenzohernandez/transactions-processor/internal/repository"
	s3_services "github.com/rodrigoenzohernandez/transactions-processor/internal/services/s3"
	sqs_services "github.com/rodrigoenzohernandez/transactions-processor/internal/services/sqs"

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
	// Get bucket and key from the event
	eventJson, _ := json.Marshal(event)
	var data Event
	json.Unmarshal(eventJson, &data)
	bucket := data.Records[0].S3.Bucket.Name
	key := data.Records[0].S3.Object.Key
	log.Info(fmt.Sprintf("Process started: Object %s uploaded to the bucket", key))

	// Validate file extension
	if !strings.HasSuffix(key, ".csv") {
		log.Error(fmt.Sprintf("File %s is not valid. It's not a .csv file", key))
		return
	}

	// Connect to the database and create a transactions repo
	db, _ := repository.Connect()
	defer repository.Disconnect(db)
	transactionsRepo := repository.NewTransactionRepo(db)

	// Define session
	session := session.Must(session.NewSession())

	// Define clients
	s3Client := s3.New(session)
	sqsClient := sqs.New(session)

	object, err := s3_services.GetObject(bucket, key, s3Client)
	if err != nil {
		return
	}

	defer object.Body.Close()

	records, err := utils.GetRecordsFromBuffer(object.Body)
	if err != nil {
		return
	}

	e := transactionsRepo.InsertMany(records)
	if e != nil {
		return
	}

	report, err := json.Marshal(utils.GenerateReport(records))
	if err != nil {
		log.Debug(fmt.Sprintf("Issue with the report: %s", err))
		return
	}
	log.Debug(fmt.Sprintf("Report: %s", report))

	queueUrl, err := sqs_services.GetQueueURL("reports_queue", sqsClient)
	if err != nil {
		return
	}

	sqs_services.SendMessage(report, queueUrl, sqsClient)

}

func main() {
	lambda.Start(handler)
}
