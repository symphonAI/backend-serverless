package main

import "github.com/aws/aws-lambda-go/events"

func generatePolicy(methodArn string) (events.APIGatewayCustomAuthorizerPolicy){
	policyDocument := events.APIGatewayCustomAuthorizerPolicy{
		Version: "2012-10-17",
		Statement: []events.IAMPolicyStatement{
			{
				Action:   []string{"execute-api:Invoke"},
				Effect:   "Allow",
				Resource: []string{methodArn},
			},
		},
	}
	return policyDocument

}