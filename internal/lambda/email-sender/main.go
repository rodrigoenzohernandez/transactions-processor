package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rodrigoenzohernandez/transactions-processor/internal/services"
	ssm_services "github.com/rodrigoenzohernandez/transactions-processor/internal/services/ssm"
	"github.com/rodrigoenzohernandez/transactions-processor/internal/types"
	"github.com/rodrigoenzohernandez/transactions-processor/internal/utils"
	"github.com/rodrigoenzohernandez/transactions-processor/internal/utils/logger"
)

var log = logger.GetLogger("email_sender_lambda")

func handler(ctx context.Context, sqsEvent events.SQSEvent) {
	for _, message := range sqsEvent.Records {
		log.Info(fmt.Sprintf("Process started. Message %s has been sent to the SQS Queue", message.MessageId))

		var data types.Report
		json.Unmarshal([]byte(message.Body), &data)

		emailContent, _ := utils.GenerateEmailContent(data)

		log.Info(fmt.Sprintf("Email content: %v", emailContent))

		notificationEmail, _ := ssm_services.GetSSMParameter("/smtp/provider/sender")

		services.SendEmail(notificationEmail, "Transactions report", emailContent)

	}

}

func main() {
	lambda.Start(handler)
}
