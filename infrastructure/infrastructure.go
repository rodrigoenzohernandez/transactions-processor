package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
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

	var env string
	env = os.Getenv("ENV")
	if env == "" {
		env = "develop"
	}

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

	paramUpdaterLambda := awslambda.NewFunction(stack, jsii.String("ssm-params-updater"), &awslambda.FunctionProps{
		FunctionName: jsii.String("ssm-params-updater"),
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.Code_FromAsset(jsii.String("../internal/lambda/ssm-params-updater"), nil),
	})

	ssmParamsRestApi := awsapigateway.NewRestApi(stack, jsii.String("ssm-params-api"), &awsapigateway.RestApiProps{
		RestApiName: jsii.String("ssmParamsApi"),
		Description: jsii.String("API to update SSM parameters"),
		DeployOptions: &awsapigateway.StageOptions{
			StageName: jsii.String(env),
		},
	})

	updateResource := ssmParamsRestApi.Root().AddResource(jsii.String("param"), nil)
	updateResource.AddMethod(jsii.String("PUT"), awsapigateway.NewLambdaIntegration(paramUpdaterLambda, nil), &awsapigateway.MethodOptions{
		AuthorizationType: awsapigateway.AuthorizationType_NONE,
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

	awsssm.NewStringParameter(stack, jsii.String("SMTP_PROVIDER_SENDER"), &awsssm.StringParameterProps{
		ParameterName: jsii.String("/smtp/provider/sender"),
		StringValue:   jsii.String(os.Getenv("SMTP_PROVIDER_SENDER")),
	})

	awsssm.NewStringParameter(stack, jsii.String("SMTP_NOTIFICATION_EMAIL"), &awsssm.StringParameterProps{
		ParameterName: jsii.String("/smtp/notification/email"),
		StringValue:   jsii.String(os.Getenv("SMTP_NOTIFICATION_EMAIL")),
	})

	// permissions

	transactionsBucket.GrantRead(filesProcessorLambda, nil)
	reportsQueue.GrantSendMessages(filesProcessorLambda)

	paramUpdaterLambda.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("ssm:PutParameter"),
		Resources: jsii.Strings(fmt.Sprintf("arn:aws:ssm:%s:%s:parameter/smtp/notification/email", *stack.Region(), *stack.Account())),
	}))

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
