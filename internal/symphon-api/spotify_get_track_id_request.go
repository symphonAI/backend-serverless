package symphonapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) getSpotifyTrackID(spotifyAccessToken string, trackName string, artistName string, trackIDChannel chan SpotifyTrackIDResult) {
	trackName = removeUnsupportedCharacters(trackName)
	artistName = removeUnsupportedCharacters(artistName)
	query := "/search?q=track:" + url.QueryEscape(trackName) + "%20artist:" + url.QueryEscape(artistName) + "&type=track&limit=1"
	endpoint := SPOTIFY_BASE_URL + query
	fmt.Println("Making GET request to:", endpoint)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		trackIDChannel <- SpotifyTrackIDResult{Error: err}
		return
	}

	req.Header.Add("Authorization", "Bearer "+spotifyAccessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		trackIDChannel <- SpotifyTrackIDResult{Error: err}
		return
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		trackIDChannel <- SpotifyTrackIDResult{Error: err}
		return
	}

	spotifyResponse := SpotifyTrackIDResponse{}
	err = json.Unmarshal(responseBody, &spotifyResponse)
	if err != nil {
		trackIDChannel <- SpotifyTrackIDResult{Error: err}
		return
	}
	if len(spotifyResponse.Tracks.Items) == 0 {
		fmt.Println("WARNING: Could not find track:", trackName, artistName)
		return
	}
	if len(spotifyResponse.Tracks.Items[0].ID) < 1 {
		trackIDChannel <- SpotifyTrackIDResult{Error: fmt.Errorf("cannot find song in Spotify")}
	} else {
		trackID := "spotify:track:" + spotifyResponse.Tracks.Items[0].ID
		fmt.Println("Track ID:", trackID)
		trackIDChannel <- SpotifyTrackIDResult{ID: trackID}
	}
}

func (c *Client) GetAllSpotifyTrackIDs(spotifyAccessToken string, chatGPTRecommendations ChatGPTRecommendations) ([]string, error) {

	trackIDChannel := make(chan SpotifyTrackIDResult, len(chatGPTRecommendations))
	defer close(trackIDChannel)

	for _, recommendation := range chatGPTRecommendations {
		go c.getSpotifyTrackID(spotifyAccessToken, recommendation.Track, recommendation.Artist, trackIDChannel)
	}

	trackIDs := []string{}
	err := error(nil)

	for range chatGPTRecommendations {
		trackIDResult := <-trackIDChannel
		if trackIDResult.Error != nil {
			err = trackIDResult.Error
			continue
		}
		trackIDs = append(trackIDs, trackIDResult.ID)
	}

	return trackIDs, err
}

func removeUnsupportedCharacters(input string) string {
	output := strings.ReplaceAll(input, "'", "")
	return output
}