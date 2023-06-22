package symphonapi

type UserFields struct {
	Prompt      string
	Temperature string
	Options     []string
}

type Track struct {
	Title string `json:"track"`
	Artist string `json:"artist"`
}

type SpotifyTrackResponse struct {
	Items []struct {
		Artists []struct {
			Name string `json:"name"`
		} `json:"artists"`
		Name string `json:"name"`
	} `json:"items"`
}

type SpotifyBandResult struct {
	Items []struct {
		Name string `json:"name"`
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

type SpotifyTrackIDResponse struct {
	Tracks struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	} `json:"tracks"`
}

type SpotifyTrackIDResult struct {
	ID    string
	Error error
}

type CreatePlaylistResponse struct {
	ID  string `json:"id"`
	URI string `json:"uri"`
}

type OpenAIModel interface {
	GeneratePayload(userFields UserFields) (map[string]interface{}, error)
	GetUrl() string
	ParseRecommendedTracksFromResponse([]byte) ([]Track, error)
}