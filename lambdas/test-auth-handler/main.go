// main.go
package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)


func testAuth(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Test auth:", request.RequestContext.Authorizer)
	accessToken := request.RequestContext.Authorizer["lambda"].(map[string]interface{})["accessToken"].(string)

	fmt.Println("Access token:", accessToken)
	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body: accessToken,
	}
	return response, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(testAuth)
}
