package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type RequestBody struct {
	NotificationEmail string `json:"notificationEmail"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body RequestBody
	err := json.Unmarshal([]byte(request.Body), &body)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("Invalid request body: %v", err),
		}, nil
	}

	session := session.Must(session.NewSession())
	ssmClient := ssm.New(session)

	paramName := "/smtp/notification/email"
	paramValue := body.NotificationEmail

	_, err = ssmClient.PutParameter(&ssm.PutParameterInput{
		Name:      aws.String(paramName),
		Value:     aws.String(paramValue),
		Overwrite: aws.Bool(true),
		Type:      aws.String("String"),
	})
	if err != nil {
		log.Printf("Error updating SSM parameter: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Failed to update parameter: %v", err),
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
