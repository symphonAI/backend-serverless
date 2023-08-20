package symphonapi

import (
	"encoding/json"
	"strconv"
	"strings"
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
	var tracks []Track

	data := "[" + strings.ReplaceAll(response.Choices[0].ChatGPTMessage.Content, "\n", "") + "]"

	err = json.Unmarshal([]byte(data), &tracks)
	if err != nil {
		return []Track{}, err
	}
	return tracks, nil
}

type ChatGPTFunctionResponse struct {
	Choices []ChatGPTFunctionChoice `json:"choices"` 
}

type ChatGPTFunctionChoice struct {
	ChatGPTMessage ChatGPTMessage `json:"message"`
}

type ChatGPTMessage struct {
	Content string `json:"content"`
}
