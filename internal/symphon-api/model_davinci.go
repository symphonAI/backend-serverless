package symphonapi

import (
	"encoding/json"
	"strconv"
)

type DaVinciModel struct {}

func (m *DaVinciModel) GetUrl() string {
	return OPENAI_BASE_URL + "/completions"
}

func (m *DaVinciModel) GeneratePayload(userFields UserFields) (map[string]interface{}, error) {
	floatTemperature, err := strconv.ParseFloat(userFields.Temperature, 64)
	if err != nil {
		return nil, err
	}
	payload := map[string]interface{}{
		"model":       "text-davinci-003",
		"max_tokens":  2000,
		"temperature": floatTemperature,
		"prompt":      userFields.Prompt,
	}
	return payload, nil
}

func (m *DaVinciModel) ParseRecommendedTracksFromResponse(responseBody []byte) ([]Track, error) {
	response := ChatGPTResponse{}
	err := json.Unmarshal(responseBody, &response)
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
