package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/awslambdago"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
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

	// DynamoDB table
	// Going with single-table design here,
	// So given generic names to the PartitionKey
	// and SortKey (literally just "PartitionKey" and "SortKey")
	// For more info, please see:
	// https://www.alexdebrie.com/posts/dynamodb-single-table/
	ddbTable := awsdynamodb.NewTable(stack, jsii.String("symphonai-dbtbl"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("PartitionKey"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: jsii.String("SortKey"),
			Type: awsdynamodb.AttributeType_STRING,
		},
	})

	// Create a new api HTTP api on gateway v2.
	max_age_in_minutes := 10.00
	api := awsapigatewayv2.NewHttpApi(stack, jsii.String("symphonai-api"), &awsapigatewayv2.HttpApiProps{
		CorsPreflight: &awsapigatewayv2.CorsPreflightOptions{
			AllowOrigins: &[]*string{jsii.String("https://symphon.ai"), jsii.String("http://localhost:3000")}, // Provide a list of allowed origins
			AllowMethods: &[]awsapigatewayv2.CorsHttpMethod{
				awsapigatewayv2.CorsHttpMethod_ANY,
			},
			AllowHeaders:  &[]*string{jsii.String("*")},
			ExposeHeaders: &[]*string{jsii.String("*")},
			MaxAge:        awscdk.Duration_Minutes(&max_age_in_minutes),
		},
	})

	customAuthorizerEnvVars := make(map[string]*string)

	addSecretCredentialsToEnvVars(customAuthorizerEnvVars)

	// Custom Authorizer
	customAuthorizerFunc := awslambdago.NewGoFunction(stack, jsii.String("test-auth-handler"), &awslambdago.GoFunctionProps{
		MemorySize:  jsii.Number(128),
		ModuleDir:   jsii.String("./go.mod"),
		Entry:       jsii.String("./lambdas/custom-authorizer"),
		Environment: &customAuthorizerEnvVars,
		Runtime: awslambda.Runtime_GO_1_X(),
	})

	str := "$request.header.Authorization"
	identitySources := []*string{}
	identitySources = append(identitySources, &str)

	authorizer := awsapigatewayv2.NewHttpAuthorizer(
		scope, 
		jsii.String("custom-authorizer"),
		&awsapigatewayv2.HttpAuthorizerProps{
				HttpApi: api,
				IdentitySource: &identitySources,
				Type: awsapigatewayv2.HttpAuthorizerType_LAMBDA,
				AuthorizerUri: customAuthorizerFunc.FunctionArn(),
		})

	// Prompt lambda function.

	promptLambdaEnvVars := make(map[string]*string)

	addSecretCredentialsToEnvVars(promptLambdaEnvVars)

	durationInMinutes := 10.00

	promptFunc := awslambdago.NewGoFunction(stack, jsii.String("prompt-handler"), &awslambdago.GoFunctionProps{
		MemorySize: jsii.Number(128),
		ModuleDir:  jsii.String("./go.mod"),
		Entry:      jsii.String("./lambdas/prompt-handler"),
		Environment: &promptLambdaEnvVars,
		Timeout:  awscdk.Duration_Minutes(&durationInMinutes),
		Runtime: awslambda.Runtime_GO_1_X(),
	})

	promptIntegration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("PromptIntegration"),
		promptFunc,
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{})

	api.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Integration: promptIntegration,
		Path:        jsii.String("/prompt"),
		Authorizer: awsapigatewayv2.HttpAuthorizer_FromHttpAuthorizerAttributes(
			scope, 
			jsii.String("custom-auth"), 
			&awsapigatewayv2.HttpAuthorizerAttributes{
				AuthorizerId: authorizer.AuthorizerId(),
				AuthorizerType: jsii.String("CUSTOM"),
			}),
		Methods: &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_POST},
	})

	// Signup lambda function

	// First, initialise the environment variables
	// that will go into the lambda function
	signupLambdaEnvVars := make(map[string]*string)

	// Assign the values to the map
	signupLambdaEnvVars["DYNAMODB_TABLE_NAME"] = ddbTable.TableName()

	addSecretCredentialsToEnvVars(signupLambdaEnvVars)

	// The role is gonna be a bit of a pain...
	signupRole := awsiam.NewRole(stack, jsii.String("signup-lambda-role"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("lambda.amazonaws.com"), &awsiam.ServicePrincipalOpts{
			Region: jsii.String("ap-southeast-2"),
		}),
	})

	// Define the policy statements
	statements := []awsiam.PolicyStatement{}

	// Enable logging Lambda function to Cloudwatch
	statements = append(statements,
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			Effect:    awsiam.Effect_ALLOW,
			Resources: jsii.Strings("*"),
			Actions: &[]*string{
				jsii.String("logs:CreateLogGroup"),
				jsii.String("logs:CreateLogStream"),
				jsii.String("logs:PutLogEvents"),
			},
		}),
	)

	ddbTableInPolicy := aws.StringSlice([]string{*ddbTable.TableArn()})
	fmt.Println("DynamoDB resource:", ddbTableInPolicy)

	statements = append(statements,
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			Effect:    awsiam.Effect_ALLOW,
			Resources: &ddbTableInPolicy,
			Actions:   jsii.Strings("*"),
		}),
	)

	signupRoleInPolicy := []awsiam.IRole{signupRole}
	fmt.Println("SignupRole resource:", signupRoleInPolicy)

	// Create the policy and add the statements and add the role
	awsiam.NewPolicy(stack, jsii.String("signup-role-policy"), &awsiam.PolicyProps{
		PolicyName: aws.String("signup-role-policy"),
		Statements: &statements,
		Roles:      &signupRoleInPolicy,
	})

	signupFunc := awslambdago.NewGoFunction(stack, jsii.String("signup-handler"), &awslambdago.GoFunctionProps{
		MemorySize:  jsii.Number(128),
		ModuleDir:   jsii.String("./go.mod"),
		Entry:       jsii.String("./lambdas/signup-handler"),
		Environment: &signupLambdaEnvVars,
		Role:        signupRole,
		Runtime: awslambda.Runtime_GO_1_X(),
	})

	signupIntegration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("SignupIntegration"),
		signupFunc,
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{})

	api.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Integration: signupIntegration,
		Path:        jsii.String("/signup"),
		Methods: &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_POST},
	})

	// Test-Auth Handler
	testAuthFunc := awslambdago.NewGoFunction(stack, jsii.String("test-auth-handler"), &awslambdago.GoFunctionProps{
		MemorySize:  jsii.Number(128),
		ModuleDir:   jsii.String("./go.mod"),
		Entry:       jsii.String("./lambdas/test-auth-handler"),
		Environment: &signupLambdaEnvVars,
		Runtime: awslambda.Runtime_GO_1_X(),
	})

	testAuthIntegration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("SignupIntegration"),
		testAuthFunc,
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{})

	api.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Integration: testAuthIntegration,
		Path:        jsii.String("/test-auth"),
		Methods: &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_GET},
		Authorizer: awsapigatewayv2.HttpAuthorizer_FromHttpAuthorizerAttributes(
			scope, 
			jsii.String("custom-auth"), 
			&awsapigatewayv2.HttpAuthorizerAttributes{
				AuthorizerId: authorizer.AuthorizerId(),
				AuthorizerType: jsii.String("CUSTOM"),
			}),
	})

	return stack
}

func addSecretCredentialsToEnvVars(envVars map[string]*string) error {
	fmt.Println("Attempting to fetch secret credentials from SSM parameter store...")
	cfg, err := config.LoadDefaultConfig(context.TODO())
	cfg.Region = "ap-southeast-2"
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	ssmClient := ssm.NewFromConfig(cfg)

	parameterName := "/symphonai/credentials/prod"

	decrypt := true
	input := &ssm.GetParameterInput{
		Name:           &parameterName,
		WithDecryption: &decrypt,
	}

	result, err := ssmClient.GetParameter(context.TODO(), input)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if err != nil {
		return err
	}

	// Parse the JSON string into a map[string]interface{}
	var credentials map[string]interface{}
	err = json.Unmarshal([]byte(*result.Parameter.Value), &credentials)
	if err != nil {
		return err
	}

	// Assign each key/value pair to the input map[string]*string parameter
	for key, value := range credentials {
		strValue := fmt.Sprintf("%v", value)
		envVars[key] = &strValue
	}

	return nil
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
