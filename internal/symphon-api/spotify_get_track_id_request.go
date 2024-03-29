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
		eMsg := fmt.Sprintf("Could not find track: %v - %v", trackName, artistName)
		fmt.Println(eMsg)
		trackIDChannel <- SpotifyTrackIDResult{Error: fmt.Errorf(eMsg)}
		return
	}
	if len(spotifyResponse.Tracks.Items[0].ID) < 1 {
		trackIDChannel <- SpotifyTrackIDResult{Error: fmt.Errorf("cannot find song in Spotify")}
	} else {
		trackID := "spotify:track:" + spotifyResponse.Tracks.Items[0].ID
		fmt.Printf("Found track: %v - %v with Spotify Track ID: %v\n", trackName, artistName, trackID)
		trackIDChannel <- SpotifyTrackIDResult{ID: trackID}
	}
}

func (c *Client) GetAllSpotifyTrackIDs(spotifyAccessToken string, recommendedTracks []Track) ([]string, error) {

	trackIDChannel := make(chan SpotifyTrackIDResult, len(recommendedTracks))
	defer close(trackIDChannel)

	for _, recommendedTrack := range recommendedTracks {
		go c.getSpotifyTrackID(spotifyAccessToken, recommendedTrack.Title, recommendedTrack.Artist, trackIDChannel)
	}

	trackIDs := []string{}
	err := error(nil)

	for range recommendedTracks {
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