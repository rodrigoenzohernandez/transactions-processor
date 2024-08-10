package ssm_services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// Put SSM param value overwriting
func PutSSMParameter(paramName string, paramValue string, dataType string) (string, error) {
	session := session.Must(session.NewSession())
	ssmClient := ssm.New(session)

	_, err := ssmClient.PutParameter(&ssm.PutParameterInput{
		Name:      aws.String(paramName),
		Value:     aws.String(paramValue),
		Overwrite: aws.Bool(true),
		Type:      aws.String(dataType),
	})

	return "", err
}
