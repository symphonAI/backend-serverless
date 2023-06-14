package symphonapi

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (c *Client) GetTopTracksSpotify(spotifyAccessToken string, trackChannel chan SpotifyResult) {

	endpoint := SPOTIFY_BASE_URL + "/me/top/tracks?limit=25&time_range=long_term"

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		trackChannel <- SpotifyResult{
			Message: nil,
			Error:   err,
		}
		return
	}

	req.Header.Add("Authorization", "Bearer "+spotifyAccessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		trackChannel <- SpotifyResult{
			Message: nil,
			Error:   err,
		}
		return
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		trackChannel <- SpotifyResult{
			Message: nil,
			Error:   err,
		}
		return
	}

	spotifyResponse := SpotifyTrackResponse{}
	err = json.Unmarshal(responseBody, &spotifyResponse)
	if err != nil {
		trackChannel <- SpotifyResult{
			Message: nil,
			Error:   err,
		}
		return
	}

	topTracks := []string{}
	for _, track := range spotifyResponse.Items {
		artist := track.Artists[0].Name
		trackName := track.Name
		topTracks = append(topTracks, trackName+" - "+artist)
	}

	trackChannel <- SpotifyResult{
		Message: topTracks,
		Error:   nil,
	}
}

// topTracks := []string{
// 	"Use Somebody - Kings of Leon",
// 	"Last Nite - The Strokes",
// 	"Mr. Brightside - The Killers",
// 	"Seven Nation Army - The White Stripes",
// 	"Lonely Boy - The Black Keys",
// 	"I Bet You Look Good On The Dancefloor - Arctic Monkeys",
// 	"Hate To Say I Told You So - The Hives",
// 	"Get Free - The Vines",
// 	"Can't Stand Me Now - The Libertines",
// 	"Steady, As She Goes - The Raconteurs",
// 	"Chelsea Dagger - The Fratellis",
// 	"If You Wanna - The Vaccines",
// 	"Naive - The Kooks",
// 	"Let's Dance To Joy Division - The Wombats",
// 	"Standing Next To Me - The Last Shadow Puppets",
// }
