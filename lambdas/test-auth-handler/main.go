// main.go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)


func testAuth(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	
	fmt.Println("Test auth:", request.RequestContext.Authorizer)
	accessToken := request.RequestContext.Authorizer["lambda"].(map[string]interface{})["accessToken"].(string)

	// Have to do this annoying workaround because  
	// SAM CLI https://github.com/aws/aws-sam-cli/issues/4161
	// TODO remove this when the issue is fixed
	headers := make(map[string]string)
	if os.Getenv("ENV") == "dev" {
		headers["Access-Control-Allow-Origin"] = "http://localhost:3000"
		headers["Access-Control-Allow-Credentials"] = "true"
	}
	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body: accessToken,
		Headers: headers,
	}
	return response, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(testAuth)
}
