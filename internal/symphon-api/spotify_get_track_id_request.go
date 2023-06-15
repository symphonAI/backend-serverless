package symphonapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
)

func (c *Client) getSpotifyTrackID(wg sync.WaitGroup, spotifyAccessToken string, trackName string, artistName string, trackIDChannel chan SpotifyTrackIDResult) {
	defer wg.Done()

	endpoint := SPOTIFY_BASE_URL + "/search?q=track:" + trackName + "%20artist:" + artistName + "&type=track&limit=1"
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
	trackID := spotifyResponse.Tracks.Items[0].ID

	trackIDChannel <- SpotifyTrackIDResult{ID: trackID}
}

func (c *Client) GetAllSpotifyTrackIDs(spotifyAccessToken string, chatGPTRecommendations ChatGPTRecommendations) ([]string, error) {
	var wg sync.WaitGroup
	wg.Add(len(chatGPTRecommendations))

	trackIDChannel := make(chan SpotifyTrackIDResult)
	defer close(trackIDChannel)

	trackIDs := []string{}
	for _, recommendation := range chatGPTRecommendations {
		go c.getSpotifyTrackID(spotifyAccessToken, recommendation.Track, recommendation.Artist, trackIDChannel)
	}

	for range trackIDs {
		trackIDResult := <-trackIDChannel
		if trackIDResult.Error != nil {
			return nil, trackIDResult.Error
		}
		trackIDs = append(trackIDs, trackIDResult.ID)
	}

	wg.Wait()

	return trackIDs, nil
}
