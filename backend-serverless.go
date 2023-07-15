package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/awscertificatemanager"
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

	env := os.Getenv("ENV")
	fmt.Println("Deploying for env:", env)

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
	allowCredentials := true
	corsOrigin := "https://symphon.ai" // URL of the website on AWS
	if env == "dev" {
		corsOrigin = "http://localhost:3000" // URL of the website when running locally
	}
	api := awsapigatewayv2.NewHttpApi(stack, jsii.String("symphonai-api"), &awsapigatewayv2.HttpApiProps{
		CorsPreflight: &awsapigatewayv2.CorsPreflightOptions{
			AllowOrigins: &[]*string{jsii.String(corsOrigin)}, // Provide a list of allowed origins
			AllowMethods: &[]awsapigatewayv2.CorsHttpMethod{
				awsapigatewayv2.CorsHttpMethod_ANY,
			},
			AllowCredentials: &allowCredentials,
			AllowHeaders:  &[]*string{jsii.String("Content-Type")},
			ExposeHeaders: &[]*string{jsii.String("*")},
			MaxAge:        awscdk.Duration_Minutes(&max_age_in_minutes),
		},
	})

	apiGwDomainName := awsapigatewayv2.NewDomainName(stack, jsii.String("symphonaiapigwdomain"), &awsapigatewayv2.DomainNameProps{
		Certificate: awscertificatemanager.Certificate_FromCertificateArn(
			stack, jsii.String("symphonaicert"), 
			jsii.String("arn:aws:acm:ap-southeast-2:349564020337:certificate/31bac708-a985-4350-b584-e21abc042cbc")),
		DomainName: jsii.String("api.symphon.ai"),
		EndpointType: awsapigatewayv2.EndpointType_REGIONAL,
	})

	awsapigatewayv2.NewApiMapping(stack, jsii.String("symphonai-api-mapping"), &awsapigatewayv2.ApiMappingProps{
		Api: api,
		DomainName: apiGwDomainName,
	})

	// This part, specifically creation of the new A Record is just hanging
	// TODO figure out why
	// symphonAiHostedZone := awsroute53.HostedZone_FromHostedZoneAttributes(stack, jsii.String("symphonai-hz"), &awsroute53.HostedZoneAttributes{
	// 	HostedZoneId: jsii.String("Z0045186G26CSUNVWCC5"),
	// 	ZoneName: jsii.String("symphon.ai"),
	// })
	// awsroute53.NewARecord(stack, jsii.String("apirecord"), &awsroute53.ARecordProps{
	// 	Zone: symphonAiHostedZone,
	// 	// RecordName: jsii.String("api.symphon.ai"),
	// 	Target: awsroute53.RecordTarget_FromAlias(
	// 		awsroute53targets.NewApiGatewayv2DomainProperties(
	// 			apiGwDomainName.RegionalDomainName(), 
	// 			apiGwDomainName.RegionalHostedZoneId()),
	// )})

	customAuthorizerEnvVars := make(map[string]*string)
	customAuthorizerEnvVars["ENV"] = &env
	customAuthorizerEnvVars["DYNAMODB_TABLE_NAME"] = ddbTable.TableName()
	addSecretCredentialsToEnvVars(customAuthorizerEnvVars)

	// Custom Authorizer - Role
	customAuthorizerRole := awsiam.NewRole(stack, jsii.String("custom-authorizer-role"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("lambda.amazonaws.com"), &awsiam.ServicePrincipalOpts{
			Region: jsii.String("ap-southeast-2"),
		}),
	})

	// Define the policy statements
	customAuthStatements := []awsiam.PolicyStatement{}

	// Enable logging Lambda function to Cloudwatch
	customAuthStatements = append(customAuthStatements,
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

	customAuthStatements = append(customAuthStatements,
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			Effect:    awsiam.Effect_ALLOW,
			Resources: &ddbTableInPolicy,
			Actions:   jsii.Strings("dynamodb:GetItem"),
		}),
	)

	customAuthRoleInPolicy := []awsiam.IRole{customAuthorizerRole}
	fmt.Println("CustomAuthRole resource:", customAuthRoleInPolicy)

	awsiam.NewPolicy(stack, jsii.String("customauth-role-policy"), &awsiam.PolicyProps{
		PolicyName: aws.String("customauth-role-policy"),
		Statements: &customAuthStatements,
		Roles:      &customAuthRoleInPolicy,
	})

	// Custom Authorizer
	customAuthorizerFunc := awslambdago.NewGoFunction(stack, jsii.String("custom-authorizer-lambda"), &awslambdago.GoFunctionProps{
		MemorySize:  jsii.Number(128),
		ModuleDir:   jsii.String("./go.mod"),
		Entry:       jsii.String("./lambdas/custom-authorizer"),
		Environment: &customAuthorizerEnvVars,
		Runtime: awslambda.Runtime_GO_1_X(),
		Role: customAuthorizerRole,
	})

	str := "$request.header.cookie"
	identitySources := []*string{}
	identitySources = append(identitySources, &str)

	authorizerLambdaArn := customAuthorizerFunc.FunctionArn()
	authorizerUri := "arn:aws:apigateway:ap-southeast-2:lambda:path/2015-03-31/functions/" + *authorizerLambdaArn + "/invocations"
	authorizer := awsapigatewayv2.NewHttpAuthorizer(
		stack, 
		jsii.String("custom-authorizer"),
		&awsapigatewayv2.HttpAuthorizerProps{
				HttpApi: api,
				Type: awsapigatewayv2.HttpAuthorizerType_LAMBDA,
				AuthorizerUri: &authorizerUri,
				IdentitySource: &identitySources,
		})


	// Prompt lambda function.
	promptLambdaEnvVars := make(map[string]*string)

	/*
		Current available models:
			"davinci"
			"gpt-3.5-turbo"
	*/
	promptLambdaEnvVars["OPENAI_MODEL"] = jsii.String("gpt-3.5-turbo")
	promptLambdaEnvVars["ENV"] = &env
	
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
			stack, 
			jsii.String("custom-auth-for-prompt"), 
			&awsapigatewayv2.HttpAuthorizerAttributes{
				AuthorizerId: authorizer.AuthorizerId(),
				AuthorizerType: jsii.String("CUSTOM"),
			}),
		Methods: &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_POST},
	})

	// Login lambda function

	// First, initialise the environment variables
	// that will go into the lambda function
	loginLambdaEnvVars := make(map[string]*string)

	// Assign the values to the map
	loginLambdaEnvVars["DYNAMODB_TABLE_NAME"] = ddbTable.TableName()
	loginLambdaEnvVars["ENV"] = &env

	addSecretCredentialsToEnvVars(loginLambdaEnvVars)

	// The role is gonna be a bit of a pain...
	loginRole := awsiam.NewRole(stack, jsii.String("login-lambda-role"), &awsiam.RoleProps{
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

	statements = append(statements,
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			Effect:    awsiam.Effect_ALLOW,
			Resources: &ddbTableInPolicy,
			Actions:   jsii.Strings("*"),
		}),
	)

	loginRoleInPolicy := []awsiam.IRole{loginRole}
	fmt.Println("loginRole resource:", loginRoleInPolicy)

	// Create the policy and add the statements and add the role
	awsiam.NewPolicy(stack, jsii.String("login-role-policy"), &awsiam.PolicyProps{
		PolicyName: aws.String("login-role-policy"),
		Statements: &statements,
		Roles:      &loginRoleInPolicy,
	})

	loginFunc := awslambdago.NewGoFunction(stack, jsii.String("login-handler"), &awslambdago.GoFunctionProps{
		MemorySize:  jsii.Number(128),
		ModuleDir:   jsii.String("./go.mod"),
		Entry:       jsii.String("./lambdas/login-handler"),
		Environment: &loginLambdaEnvVars,
		Role:        loginRole,
		Runtime: awslambda.Runtime_GO_1_X(),
	})

	loginIntegration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("loginIntegration"),
		loginFunc,
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{})

	api.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Integration: loginIntegration,
		Path:        jsii.String("/login"),
		Methods: &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_POST},
	})

	// Test-Auth Handler
	testAuthFunc := awslambdago.NewGoFunction(stack, jsii.String("test-auth-handler"), &awslambdago.GoFunctionProps{
		MemorySize:  jsii.Number(128),
		ModuleDir:   jsii.String("./go.mod"),
		Entry:       jsii.String("./lambdas/test-auth-handler"),
		Environment: &loginLambdaEnvVars,
		Runtime: awslambda.Runtime_GO_1_X(),
	})

	testAuthIntegration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("loginIntegration"),
		testAuthFunc,
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{})

	api.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Integration: testAuthIntegration,
		Path:        jsii.String("/test-auth"),
		Methods: &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_GET},
		Authorizer: awsapigatewayv2.HttpAuthorizer_FromHttpAuthorizerAttributes(
			stack, 
			jsii.String("custom-auth-for-test"), 
			&awsapigatewayv2.HttpAuthorizerAttributes{
				AuthorizerId: authorizer.AuthorizerId(),
				AuthorizerType: jsii.String("CUSTOM"),
			}),
	})

	// Logout Handler
	logoutFunc := awslambdago.NewGoFunction(stack, jsii.String("logout-handler"), &awslambdago.GoFunctionProps{
		MemorySize:  jsii.Number(128),
		ModuleDir:   jsii.String("./go.mod"),
		Entry:       jsii.String("./lambdas/logout-handler"),
		Environment: &loginLambdaEnvVars,
		Runtime: awslambda.Runtime_GO_1_X(),
	})
	
	logoutIntegration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("LogoutIntegration"),
		logoutFunc,
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{})
	
	api.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Integration: logoutIntegration,
		Path:        jsii.String("/logout"),
		Methods: &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_GET},
		Authorizer: awsapigatewayv2.HttpAuthorizer_FromHttpAuthorizerAttributes(
			stack, 
			jsii.String("custom-auth-for-logout"), 
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
