package symphonapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (c *Client) PromptChatGPT(userFields UserFields) ([]Track, error) {
	endpoint := OPENAI_BASE_URL + "/completions"

	floatTemperature, err := strconv.ParseFloat(userFields.Temperature, 64)
	if err != nil {
		return []Track{}, err
	}

	payload := map[string]interface{}{
		"model":       "text-davinci-003",
		"max_tokens":  2000,
		"temperature": floatTemperature,
		"prompt":      userFields.Prompt,
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

	response := ChatGPTResponse{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return []Track{}, err
	}

	var tracks []Track
	err = json.Unmarshal([]byte(response.Choices[0].Text), &tracks)
	if err != nil {
		return []Track{}, err
	}

	return tracks, nil
}
