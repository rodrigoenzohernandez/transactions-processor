package ssm_services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/rodrigoenzohernandez/transactions-processor/internal/utils/logger"
)

var log = logger.GetLogger("ssm_services")

// Get SSM param value
func GetSSMParameter(name string) (string, error) {
	session := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
	}))

	ssmsvc := ssm.New(session)
	param, err := ssmsvc.GetParameter(&ssm.GetParameterInput{
		Name: aws.String(name),
	})
	if err != nil {
		log.Error(fmt.Sprintf("Error getting the SSM param %v", err))
		return "", err
	}
	return *param.Parameter.Value, nil
}
