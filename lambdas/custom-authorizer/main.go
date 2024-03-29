// main.go
package main

import (
	"context"
	"errors"
	"fmt"

	utils "backend-serverless/internal/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func authorize(ctx context.Context, event events.APIGatewayV2CustomAuthorizerV2Request) (events.APIGatewayCustomAuthorizerResponse, error) {
	fmt.Println("Custom authorizer called:", event)

	headers := utils.LowercaseKeyMap(event.Headers)

	cookieStr := headers["cookie"]

	tokenString := extractJwtFromCookie(cookieStr)

	if tokenString == nil {
		fmt.Println("Missing JWT from cookie")
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized") // Return a 401 Unauthorized response
	}
	jwtClaims, err := validateJWT(*tokenString)
	if err != nil {
		fmt.Println("JWT is not valid")
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized") // Return a 401 Unauthorized response
	}
	fmt.Println("JWT claims:", jwtClaims)

	user, err := getUserFromDB(jwtClaims["user"].(string))
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	}

	ip, err := getIdentityProvider(user)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	}
	accessCredentials, err := ip.getAccessCredentials(user)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	}
	// Generate the policy document for the user
	// TODO I think there will be some BS here to deal with
	policyDocument := generatePolicy(event.RouteArn)

	// Generate the authorizer response
	response := events.APIGatewayCustomAuthorizerResponse{
		PrincipalID:  user.SortKey,
		PolicyDocument: policyDocument,
		Context: accessCredentials,
	}

	return response, nil
}


func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(authorize)
}
