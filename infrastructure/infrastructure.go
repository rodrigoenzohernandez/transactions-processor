package main

import (
	"fmt"

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
	"github.com/rodrigoenzohernandez/transactions-processor/internal/utils"
)

type InfrastructureStackProps struct {
	awscdk.StackProps
}

func createQueue(stack awscdk.Stack, queueName string, visibilityTimeoutSeconds float64) awssqs.Queue {
	return awssqs.NewQueue(stack, jsii.String(queueName), &awssqs.QueueProps{
		QueueName:         jsii.String(queueName),
		VisibilityTimeout: awscdk.Duration_Seconds(jsii.Number(visibilityTimeoutSeconds)),
	})
}

func createBucket(stack awscdk.Stack, bucketName string, versioned bool) awss3.Bucket {
	return awss3.NewBucket(stack, jsii.String(bucketName), &awss3.BucketProps{
		Versioned: jsii.Bool(versioned),
	})
}

func createLambdaFunction(stack awscdk.Stack, functionName string) awslambda.Function {
	codePath := "../internal/lambda/" + functionName

	return awslambda.NewFunction(stack, jsii.String(functionName), &awslambda.FunctionProps{
		FunctionName: jsii.String(functionName),
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.Code_FromAsset(jsii.String(codePath), nil),
	})
}

func createRestApi(stack awscdk.Stack, apiName, description, stageName string) awsapigateway.RestApi {
	return awsapigateway.NewRestApi(stack, jsii.String(apiName), &awsapigateway.RestApiProps{
		RestApiName: jsii.String(apiName),
		Description: jsii.String(description),
		DeployOptions: &awsapigateway.StageOptions{
			StageName: jsii.String(stageName),
		},
	})
}

func createApiEndpoint(api awsapigateway.RestApi, resourcePath, httpMethod string, lambda awslambda.Function, authorizationType awsapigateway.AuthorizationType) {
	resource := api.Root().AddResource(jsii.String(resourcePath), nil)
	resource.AddMethod(jsii.String(httpMethod), awsapigateway.NewLambdaIntegration(lambda, nil), &awsapigateway.MethodOptions{
		AuthorizationType: authorizationType,
	})
}

func createSSMParameter(stack awscdk.Stack, parameterName, envVarName string, fallBack string) {
	awsssm.NewStringParameter(stack, jsii.String(parameterName), &awsssm.StringParameterProps{
		ParameterName: jsii.String(parameterName),
		StringValue:   jsii.String(utils.GetEnv(envVarName, fallBack)),
	})
}

func addPolicyToLambda(lambdaFunction awslambda.Function, actions []string, resources []string) {
	lambdaFunction.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings(actions...),
		Resources: jsii.Strings(resources...),
	}))
}

func NewInfrastructureStack(scope constructs.Construct, id string, props *InfrastructureStackProps) awscdk.Stack {

	env := utils.GetEnv("ENV", "develop")

	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// resources creation

	reportsQueue := createQueue(stack, "reports_queue", 300)

	transactionsBucket := createBucket(stack, "transactions-bucket", true)

	filesProcessorLambda := createLambdaFunction(stack, "files-processor")
	emailSenderLambda := createLambdaFunction(stack, "email-sender")
	paramUpdaterLambda := createLambdaFunction(stack, "ssm-params-updater")

	ssmParamsRestApi := createRestApi(stack, "ssm-params-api", "API to update SSM parameters", env)
	createApiEndpoint(ssmParamsRestApi, "param", "PUT", paramUpdaterLambda, awsapigateway.AuthorizationType_NONE)

	// SSM Parameters
	createSSMParameter(stack, "/smtp/provider/public_key", "SMTP_PROVIDER_PUBLIC_KEY", "privateKey")
	createSSMParameter(stack, "/smtp/provider/private_key", "SMTP_PROVIDER_PRIVATE_KEY", "publicKey")
	createSSMParameter(stack, "/smtp/provider/sender", "SMTP_PROVIDER_SENDER", "rodrigoenzohernandez@gmail.com")
	createSSMParameter(stack, "/smtp/notification/email", "SMTP_NOTIFICATION_EMAIL", "rodrigoenzohernandez@gmail.com")

	// permissions

	transactionsBucket.GrantRead(filesProcessorLambda, nil)
	reportsQueue.GrantSendMessages(filesProcessorLambda)

	addPolicyToLambda(paramUpdaterLambda, []string{"ssm:PutParameter"}, []string{
		fmt.Sprintf("arn:aws:ssm:%s:%s:parameter/smtp/notification/email", *stack.Region(), *stack.Account()),
	})

	addPolicyToLambda(emailSenderLambda, []string{"ssm:GetParameter"}, []string{
		fmt.Sprintf("arn:aws:ssm:%s:%s:parameter/smtp/provider/public_key", *stack.Region(), *stack.Account()),
		fmt.Sprintf("arn:aws:ssm:%s:%s:parameter/smtp/provider/private_key", *stack.Region(), *stack.Account()),
		fmt.Sprintf("arn:aws:ssm:%s:%s:parameter/smtp/provider/sender", *stack.Region(), *stack.Account()),
		fmt.Sprintf("arn:aws:ssm:%s:%s:parameter/smtp/notification/email", *stack.Region(), *stack.Account()),
	})

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
