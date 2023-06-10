package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/awscognito"
	"github.com/aws/aws-cdk-go/awscdk/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
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

	strSlice := []string{
		"email",
	}

	// Convert []string to []*string
	var ptrSlice []*string
	for _, str := range strSlice {
		ptrSlice = append(ptrSlice, &str)
	}

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

	userPoolClient := awscognito.NewCfnUserPoolClient(stack, jsii.String("SymphonAIUserPoolClient"), &awscognito.CfnUserPoolClientProps{
		UserPoolId: userPool.Ref(),
	})
	
	// Create a new api HTTP api on gateway v2.
	api := awsapigatewayv2.NewHttpApi(stack, jsii.String("symphonai-api"), &awsapigatewayv2.HttpApiProps{
		CorsPreflight: &awsapigatewayv2.CorsPreflightOptions{
			AllowOrigins: &[]*string{jsii.String("https://symphon.ai"), jsii.String("http://localhost:3000")}, // Provide a list of allowed origins	
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
	
	// First, initialise the environment variables
	// that will go into the lambda function
    signupLambdaEnvVars := make(map[string]*string)

    // Assign the values to the map
    signupLambdaEnvVars["COGNITO_USER_POOL_ID"] = userPool.Node().Id()
    signupLambdaEnvVars["COGNITO_USER_POOL_CLIENT_ID"] = userPoolClient.Node().Id()
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
	userPoolResourceInPolicy := aws.StringSlice([]string{*userPool.AttrArn()})
	fmt.Println("User Pool reference ARN:", userPoolResourceInPolicy)
	statements = append(statements, 
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			Effect:    awsiam.Effect_ALLOW,
			Resources: &userPoolResourceInPolicy,
			Actions:   jsii.Strings("*"),
		}),
	)
	// userPoolClientResourceInPolicy := aws.StringSlice([]string{*userPoolClient.Ref()})
	// fmt.Println("User Pool Client Resource:", userPoolClientResourceInPolicy)

	// statements = append(statements, 
	// 	awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
	// 		Effect:    awsiam.Effect_ALLOW,
	// 		Resources: &userPoolClientResourceInPolicy,
	// 		Actions:   jsii.Strings("*"),
	// 	}),
	// )

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
		PolicyName:     aws.String("signup-role-policy"),
		Statements:     &statements,
		Roles:          &signupRoleInPolicy,
	})

	signupFunc := awslambdago.NewGoFunction(stack, jsii.String("signup-handler"), &awslambdago.GoFunctionProps{
		MemorySize: jsii.Number(128),
		ModuleDir: jsii.String("./go.mod"),
		Entry:      jsii.String("./lambdas/signup-handler"),
		Environment: &signupLambdaEnvVars,
		Role: signupRole,
	})

	signupIntegration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("SignupIntegration"), 
		signupFunc, 
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{})

	api.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Integration: signupIntegration,
		Path:        jsii.String("/signup"),
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
		Name: &parameterName,
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
