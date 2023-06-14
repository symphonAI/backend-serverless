package symphonapi

import (
	"io/ioutil"
	"net/http"
	"encoding/json"
)

func (c *Client) GetTopBandsSpotify(spotifyAccessToken string, bandChannel chan SpotifyResult) {

	endpoint := SPOTIFY_BASE_URL + "/me/top/artists?limit=25&time_range=long_term"

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		bandChannel <- SpotifyResult{
			Message: nil,
			Error:   err,
		}
		return
	}

	req.Header.Add("Authorization", "Bearer "+spotifyAccessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		bandChannel <- SpotifyResult{
			Message: nil,
			Error:   err,
		}
		return
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		bandChannel <- SpotifyResult{
			Message: nil,
			Error:   err,
		}
		return
	}

	spotifyResponse := SpotifyBandResult{}
	err = json.Unmarshal(responseBody, &spotifyResponse)
	if err != nil {
		bandChannel <- SpotifyResult{
			Message: nil,
			Error:   err,
		}
		return
	}

	topBands := []string{}
	for _, band := range spotifyResponse.Items {
		bandName := band.Name
		topBands = append(topBands, bandName)
	}


	bandChannel <- SpotifyResult{
		Message: topBands,
		Error:   nil,
	}
}

// topBands := []string{
// 		"Kings of Leon",
// 		"The Strokes",
// 		"The Killers",
// 		"The White Stripes",
// 		"The Black Keys",
// 		"Arctic Monkeys",
// 		"The Hives",
// 		"The Vines",
// 		"The Libertines",
// 		"The Raconteurs",
// 		"The Fratellis",
// 		"The Vaccines",
// 		"The Kooks",
// 		"The Wombats",
// 		"The Last Shadow Puppets",
// 	}
