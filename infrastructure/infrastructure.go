package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambdaeventsources"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3notifications"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsssm"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type InfrastructureStackProps struct {
	awscdk.StackProps
}

func NewInfrastructureStack(scope constructs.Construct, id string, props *InfrastructureStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// resources creation

	reportsQueue := awssqs.NewQueue(stack, jsii.String("reports_queue"), &awssqs.QueueProps{
		QueueName:         jsii.String("reports_queue"),
		VisibilityTimeout: awscdk.Duration_Seconds(jsii.Number(300)),
	})

	transactionsBucket := awss3.NewBucket(stack, jsii.String("transactions-bucket"), &awss3.BucketProps{
		Versioned: jsii.Bool(true),
	})

	filesProcessorLambda := awslambda.NewFunction(stack, jsii.String("files-processor"), &awslambda.FunctionProps{
		FunctionName: jsii.String("files-processor"),
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.Code_FromAsset(jsii.String("../internal/lambda/files-processor"), nil),
	})

	emailSenderLambda := awslambda.NewFunction(stack, jsii.String("email-sender"), &awslambda.FunctionProps{
		FunctionName: jsii.String("email-sender"),
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.Code_FromAsset(jsii.String("../internal/lambda/email-sender"), nil),
	})

	// SSM Parameters
	awsssm.NewStringParameter(stack, jsii.String("SMTP_PROVIDER_PUBLIC_KEY"), &awsssm.StringParameterProps{
		ParameterName: jsii.String("/smtp/provider/public_key"),
		StringValue:   jsii.String(os.Getenv("SMTP_PROVIDER_PUBLIC_KEY")),
	})

	awsssm.NewStringParameter(stack, jsii.String("SMTP_PROVIDER_PRIVATE_KEY"), &awsssm.StringParameterProps{
		ParameterName: jsii.String("/smtp/provider/private_key"),
		StringValue:   jsii.String(os.Getenv("SMTP_PROVIDER_PRIVATE_KEY")),
	})

	awsssm.NewStringParameter(stack, jsii.String("SMTP_EMAIL_SENDER"), &awsssm.StringParameterProps{
		ParameterName: jsii.String("/smtp/email/sender"),
		StringValue:   jsii.String(os.Getenv("SMTP_EMAIL_SENDER")),
	})

	awsssm.NewStringParameter(stack, jsii.String("SMTP_EMAIL_RECEIVER"), &awsssm.StringParameterProps{
		ParameterName: jsii.String("/smtp/email/receiver"),
		StringValue:   jsii.String(os.Getenv("SMTP_EMAIL_RECEIVER")),
	})

	// permissions

	transactionsBucket.GrantRead(filesProcessorLambda, nil)
	reportsQueue.GrantSendMessages(filesProcessorLambda)

	// triggers

	transactionsBucket.AddEventNotification(awss3.EventType_OBJECT_CREATED, awss3notifications.NewLambdaDestination(filesProcessorLambda))

	emailSenderLambda.AddEventSource(awslambdaeventsources.NewSqsEventSource(reportsQueue, &awslambdaeventsources.SqsEventSourceProps{
		BatchSize: jsii.Number(1),
	}))

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewInfrastructureStack(app, "InfrastructureStack", &InfrastructureStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return nil
}
