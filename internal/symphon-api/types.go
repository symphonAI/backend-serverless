package symphonapi

type UserFields struct {
	Prompt      string
	Temperature string
	Options     []string
}

type ChatGPTResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

type SpotifyTrackResponse struct {
	Items []struct {
		Artists []struct {
			Name string `json:"name"`
		} `json:"artists"`
		Name        string `json:"name"`
	} `json:"items"`
}

type SpotifyBandResult struct {
	Items []struct {
		Name       string `json:"name"`
	} `json:"items"`
}

type SpotifyResult struct {
	Message []string
	Error   error
}


type SpotifyUserData struct {
	DisplayName string `json:"display_name"`
	ID          string `json:"id"`
	Email       string `json:"email"`
	ImageURL    string `json:"image_url"`
}
