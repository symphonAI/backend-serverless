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

	// can probably get rid of this call if we already have access to the userID somewhere else
	userData, err := cfg.symphonapiClient.GetUserDataFromSpotify(spotifyAccessToken)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	userID := userData.ID

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

	chatGPTRecommendations := symphonapi.ChatGPTRecommendations{}
	err = json.Unmarshal([]byte(chatgptResponse), &chatGPTRecommendations)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	fmt.Println("Chat GPT Recommendations:", chatGPTRecommendations)

	fmt.Println("Getting Track IDs...")
	trackIDs, err := cfg.symphonapiClient.GetAllSpotifyTrackIDs(spotifyAccessToken, chatGPTRecommendations)
	if err != nil {
		fmt.Println("Error getting track IDs:", err.Error())
	}
	fmt.Println("Received Track IDs:", trackIDs)

	fmt.Println("Creating Playlist...")
	playlistURI, err := cfg.symphonapiClient.CreateRecommendationPlaylist(spotifyAccessToken, userID, trackIDs, prompt, options)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	fmt.Println("Playlist URI:", playlistURI)

	payload := map[string]interface{}{
		"requestURI": playlistURI,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	stringBody := string(jsonPayload)

	fmt.Println("Returning response:", stringBody)

	// return response
	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       stringBody,
	}

	return response, nil
}

func main() {

	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handlePrompt)
}
