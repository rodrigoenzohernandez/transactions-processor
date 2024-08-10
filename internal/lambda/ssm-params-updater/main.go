package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ssm_services "github.com/rodrigoenzohernandez/transactions-processor/internal/services/ssm"
	"github.com/rodrigoenzohernandez/transactions-processor/internal/utils/logger"
)

type RequestBody struct {
	NotificationEmail string `json:"notificationEmail"`
}

var log = logger.GetLogger("ssm-params-updater")

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body RequestBody
	err := json.Unmarshal([]byte(request.Body), &body)
	if err != nil {
		data := fmt.Sprintf("Invalid request body: %v", err)
		log.Error(data)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       data,
		}, nil
	}

	_, err = ssm_services.PutSSMParameter("/smtp/notification/email", body.NotificationEmail, "String")

	if err != nil {
		data := fmt.Sprintf("Error updating SSM parameter: %v", err)
		log.Error(data)

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       data,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Parameter updated successfully",
	}, nil
}

func main() {
	lambda.Start(handler)
}
