// main.go
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func handlePrompt() (string, error) {
	return "Spotify playlist URL goes here!", nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handlePrompt)
}
