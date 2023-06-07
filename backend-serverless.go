package main

import (
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/awscognito"
	"github.com/aws/aws-cdk-go/awscdk/awsdynamodb"
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


	strSlice := []string{
		"email",
	}

	// Convert []string to []*string
	var ptrSlice []*string
	for _, str := range strSlice {
		ptrSlice = append(ptrSlice, &str)
	}


	// Create a AWS Cognito user pool
	userPool := awscognito.NewCfnUserPool(stack, jsii.String("symphonai-user-pool"), &awscognito.CfnUserPoolProps{
		UserPoolName: jsii.String("symphonai-user-pool"),
		AutoVerifiedAttributes: &ptrSlice,
		Schema: []interface{}{
			map[string]interface{}{
				"attributeDataType": "String",
				"name":              "email",
				"mutable":           true,
				"required":          true,
			},
			map[string]interface{}{
				"attributeDataType": "String",
				"name":              "name",
				"mutable":           true,
				"required":          false,
			},
		}})

	awscognito.NewCfnUserPoolClient(stack, jsii.String("SymphonAIUserPoolClient"), &awscognito.CfnUserPoolClientProps{
		UserPoolId: userPool.Ref(),
	})
	
	// Create a new api HTTP api on gateway v2.
	api := awsapigatewayv2.NewHttpApi(stack, jsii.String("symphonai-api"), &awsapigatewayv2.HttpApiProps{
		CorsPreflight: &awsapigatewayv2.CorsPreflightOptions{
			AllowOrigins: &[]*string{jsii.String("https://symphon.ai")}, // Provide a list of allowed origins			AllowMethods: &[]awsapigatewayv2.CorsHttpMethod{awsapigatewayv2.CorsHttpMethod_ANY},
		},
	})

	// Prompt lambda function.
	promptFunc := awslambdago.NewGoFunction(stack, jsii.String("prompt-handler"), &awslambdago.GoFunctionProps{
		MemorySize: jsii.Number(128),
		ModuleDir: jsii.String("./go.mod"),
		Entry:      jsii.String("./lambdas/prompt-handler"),
	})

	promptIntegration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("PromptIntegration"), 
		promptFunc, 
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{})

	api.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Integration: promptIntegration,
		Path:        jsii.String("/prompt"),
	})

	// Signup lambda function
	signupFunc := awslambdago.NewGoFunction(stack, jsii.String("prompt-handler"), &awslambdago.GoFunctionProps{
		MemorySize: jsii.Number(128),
		ModuleDir: jsii.String("./go.mod"),
		Entry:      jsii.String("./lambdas/signup-handler"),
	})

	signupIntegration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("SignupIntegration"), 
		signupFunc, 
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{})

	api.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Integration: signupIntegration,
		Path:        jsii.String("/signup"),
	})

	// DynamoDB table
	awsdynamodb.NewTable(stack, jsii.String("prompt-handler"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("PartitionKey"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: jsii.String("SortKey"),
			Type: awsdynamodb.AttributeType_STRING,
		},
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
