package main

import (
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/awslambdago"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)


type BackendServerlessStackProps struct {
	awscdk.StackProps
}

func NewBackendServerlessStack(scope constructs.Construct, id string, props *BackendServerlessStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// Create a new api HTTP api on gateway v2.
	api := awsapigatewayv2.NewHttpApi(stack, jsii.String("symphonai-api"), &awsapigatewayv2.HttpApiProps{
		CorsPreflight: &awsapigatewayv2.CorsPreflightOptions{
			AllowOrigins: &[]*string{jsii.String("symphon.ai")}, //
			AllowMethods: &[]awsapigatewayv2.CorsHttpMethod{awsapigatewayv2.CorsHttpMethod_ANY},
		},
	})

	// Create a new lambda function.
	promptFunc := awslambdago.NewGoFunction(stack, jsii.String("prompt-handler"), &awslambdago.GoFunctionProps{
		MemorySize: jsii.Number(128),
		Entry:      jsii.String("../lambdas/prompt-handler"),
	})

	// Add a lambda proxy integration.
	promptIntegration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("PromptIntegration"), 
		promptFunc, 
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{})

	// Add a route to api.
	api.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Integration: promptIntegration,
		Path:        jsii.String("/prompt"),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewBackendServerlessStack(app, "BackendServerlessStack", &BackendServerlessStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	return &awscdk.Environment{
	 Account: jsii.String("349564020337"),
	 Region:  jsii.String("ap-southeast-2"),
	}
}
