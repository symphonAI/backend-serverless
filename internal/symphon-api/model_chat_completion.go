package symphonapi

import (
	"encoding/json"
	"strconv"
)

type ChatCompletionModel struct {
	chatGptApiLLMModel string
}

func (m *ChatCompletionModel) GetUrl() string {
	return OPENAI_BASE_URL + "/chat/completions"
}

func (m *ChatCompletionModel) GeneratePayload(userFields UserFields) (map[string]interface{}, error) {
	floatTemperature, err := strconv.ParseFloat(userFields.Temperature, 64)
	if err != nil {
		return nil, err
	}
	
	// TODO one day this will be cleaned up...
	// but today is not that day
	payload := map[string]interface{}{
		"model":       m.chatGptApiLLMModel,
		"max_tokens":  2000,
		"temperature": floatTemperature,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": userFields.Prompt,
			},
		},
	}
	return payload, nil
}

func (m *ChatCompletionModel) ParseRecommendedTracksFromResponse(responseBody []byte) ([]Track, error) {
	response := ChatGPTFunctionResponse{}
	err := json.Unmarshal(responseBody, &response)
	if err != nil {
		return []Track{}, err
	}
	var tracklist TracklistResponse
	err = json.Unmarshal([]byte(response.Choices[0].ChatGPTFunctionMessage.ChatGPTFunctionCall.Arguments), &tracklist)
	if err != nil {
		return []Track{}, err
	}
	return tracklist.Tracks, nil
}

type ChatGPTFunctionResponse struct {
	Choices []ChatGPTFunctionChoice `json:"choices"` 
}

type ChatGPTFunctionChoice struct {
	ChatGPTFunctionMessage ChatGPTFunctionMessage `json:"message"`
}

type ChatGPTFunctionMessage struct {
	ChatGPTFunctionCall ChatGPTFunctionCall `json:"function_call"`
}

type ChatGPTFunctionCall struct {
	Name string `json:"name"`
	Arguments string `json:"arguments"`
}

type TracklistResponse struct {
	Tracks []Track `json:"tracklist"`
}