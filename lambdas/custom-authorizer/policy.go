package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func generatePolicy(methodArn string) (events.APIGatewayCustomAuthorizerPolicy){
	fmt.Println("Generating policy...")
	policyDocument := events.APIGatewayCustomAuthorizerPolicy{
		Version: "2012-10-17",
		Statement: []events.IAMPolicyStatement{
			{
				Action:   []string{"execute-api:Invoke"},
				Effect:   "Allow",
				Resource: []string{getAccessibleAPIMethodsPattern(methodArn)},
			},
		},
	}

	fmt.Println("Policy document:", policyDocument)

	return policyDocument
}

/*
	Example input: arn:aws:execute-api:ap-southeast-2:349564020337:l5gbu4y8b5/$default/GET/test-auth
	Example output: arn:aws:execute-api:ap-southeast-2:349564020337:l5gbu4y8b5/$default\/*\/* (without backslashes)
*/
func getAccessibleAPIMethodsPattern(methodArn string) string {
	// The format of ARNs is weirdly different when running locally,
	// so just allow everything when running locally
	if os.Getenv("ENV") == "dev" {
		return "*"
	}
	parts := strings.Split(methodArn, "/")
	parts[len(parts)-2] = "*"
	parts[len(parts)-1] = "*"
	return strings.Join(parts, "/")
}