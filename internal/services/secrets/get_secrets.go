package secrets_service

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/rodrigoenzohernandez/transactions-processor/internal/utils/logger"
)

var log = logger.GetLogger("secrets_service")

// Gets secret from AWS Secrets Manager and returns it as string
func GetSecret(secretName string, region string) (string, error) {

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Error(fmt.Sprintf("Failed on AWS Secrets Manager Setup: %v", err))
		return "", err
	}

	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to get the secrets from AWS Secrets Manager: %v", err))
		return "", err
	}

	var secretString string = *result.SecretString

	return secretString, nil

}
