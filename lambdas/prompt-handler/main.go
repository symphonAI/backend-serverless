// main.go
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	chatgptapi "backend-serverless/internal/chatgpt-api"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type config struct {
	chatgptClient chatgptapi.Client
}

func handlePrompt(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Prompt called")
	cfg := config{
		chatgptClient: chatgptapi.NewClient(),
	}

	requestBody := PromptRequestBody{}
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		errorString := fmt.Sprintf("unable to unmarshal request body: %s", err.Error())
		response := events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       errorString,
		}
		return response, errors.New(errorString)
	}

	fmt.Println("Received prompt:", requestBody.Prompt)
	fmt.Println("Received temperature:", requestBody.Temperature)


	// call chatGPT api
	userFields := chatgptapi.UserFields{
		Prompt:      requestBody.Prompt,
		Temperature: requestBody.Temperature,
	}
	

	chatgptResponse, err := cfg.chatgptClient.PromptChatGPT(userFields)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	fmt.Println("Chat GPT Response:", chatgptResponse)

	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       chatgptResponse,
	}

	return response, nil
}

func main() {

	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handlePrompt)
}
