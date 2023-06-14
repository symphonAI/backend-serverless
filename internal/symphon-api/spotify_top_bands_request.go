package symphonapi

import "context"

func (c *Client) GetTopBandsSpotify(ctx context.Context, bandChannel chan SpotifyResult) {
	
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
