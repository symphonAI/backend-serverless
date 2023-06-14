// main.go
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	promptengineering "backend-serverless/internal/prompt-engineering"
	symphonapi "backend-serverless/internal/symphon-api"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type config struct {
	symphonapiClient symphonapi.Client
}

func handlePrompt(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("/prompt called")
	cfg := config{
		symphonapiClient: symphonapi.NewClient(),
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

	prompt := requestBody.Prompt
	temperature := requestBody.Temperature
	options := requestBody.Options

	fmt.Println("Received prompt:", prompt)
	fmt.Println("Received temperature:", temperature)
	fmt.Println("Received options:", options)

	// access spotify token from context

	spotifyAccessToken := request.RequestContext.Authorizer["lambda"].(map[string]interface{})["accessToken"].(string)

	// get top bands and tracks concurrently
	bandChannel := make(chan symphonapi.SpotifyResult)
	trackChannel := make(chan symphonapi.SpotifyResult)

	go cfg.symphonapiClient.GetTopBandsSpotify(spotifyAccessToken, bandChannel)
	go cfg.symphonapiClient.GetTopTracksSpotify(spotifyAccessToken, trackChannel)

	topBands := <-bandChannel
	topTracks := <-trackChannel

	if topBands.Error != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	if topTracks.Error != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	fmt.Println("Top Bands:", topBands)
	fmt.Println("Top Tracks:", topTracks)

	// engineer prompt
	engineeredPrompt, err := promptengineering.EngineerPrompt(
		prompt,
		topBands.Message,
		topTracks.Message,
		options,
	)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	// call chatGPT api
	userFields := symphonapi.UserFields{
		Prompt:      engineeredPrompt,
		Temperature: temperature,
	}

	chatgptResponse, err := cfg.symphonapiClient.PromptChatGPT(userFields)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	fmt.Println("Chat GPT Response:", chatgptResponse)

	// build spotify playlist here

	// return response
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
