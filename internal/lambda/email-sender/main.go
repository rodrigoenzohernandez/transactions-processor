package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context) (string, error) {
	return "Hello from the Email Sender Lambda!", nil
}

func main() {
	lambda.Start(handler)
}