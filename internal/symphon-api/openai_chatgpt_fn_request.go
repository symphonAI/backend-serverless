package symphonapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (c *Client) PromptChatGPT(userFields UserFields) ([]Track, error) {
	endpoint := OPENAI_BASE_URL + "/chat/completions"

	floatTemperature, err := strconv.ParseFloat(userFields.Temperature, 64)
	if err != nil {
		return []Track{}, err
	}
	
	// TODO one day this will be cleaned up...
	// but today is not that day
	payload := map[string]interface{}{
		"model":       "gpt-3.5-turbo",
		"max_tokens":  2000,
		"temperature": floatTemperature,
		"messages": []map[string]interface{}{
			map[string]interface{}{
				"role":    "system",
				"content": "You are a music recommendation engine.",
			},
			map[string]interface{}{
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
	

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return []Track{}, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return []Track{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return []Track{}, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Track{}, err
	}

	response := ChatGPTFunctionResponse{}
	err = json.Unmarshal(responseBody, &response)
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
