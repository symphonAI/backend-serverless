package symphonapi

import (
	"encoding/json"
	"strconv"
)

type ChatGPT3Point5TurboModel struct {}

func (m *ChatGPT3Point5TurboModel) GetUrl() string {
	return OPENAI_BASE_URL + "/chat/completions"
}

func (m *ChatGPT3Point5TurboModel) GeneratePayload(userFields UserFields) (map[string]interface{}, error) {
	floatTemperature, err := strconv.ParseFloat(userFields.Temperature, 64)
	if err != nil {
		return nil, err
	}
	
	// TODO one day this will be cleaned up...
	// but today is not that day
	payload := map[string]interface{}{
		"model":       "gpt-3.5-turbo",
		"max_tokens":  2000,
		"temperature": floatTemperature,
		"messages": []map[string]interface{}{
			{
				"role":    "system",
				"content": "You are a music recommendation engine.",
			},
			{
				"role":    "user",
				"content": userFields.Prompt,
			},
		},
		"functions": []map[string]interface{}{
			{
				"name":        "get_tracklist",
				"description": "Gets a list of tracks from a query.",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"tracklist": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"artist": map[string]interface{}{
										"type":        "string",
										"description": "The artist of the track.",
									},
									"title": map[string]interface{}{
										"type":        "string",
										"description": "The title of the track.",
									},
								},
							},
						},
					},
					"required": []string{
						"tracklist",
					},
				},
			},
		},
		"function_call": "auto",
	}
	return payload, nil
}

func (m *ChatGPT3Point5TurboModel) ParseRecommendedTracksFromResponse(responseBody []byte) ([]Track, error) {
	response := ChatGPTFunctionResponse{}
	err := json.Unmarshal(responseBody, &response)
	if err != nil {
		return []Track{}, err
	}

	var tracks []Track
	err = json.Unmarshal([]byte(response.Choices[0].ChatGPTFunctionMessage.ChatGPTFunctionCall.Arguments), &tracks)
	if err != nil {
		return []Track{}, err
	}
	return tracks, nil
}
