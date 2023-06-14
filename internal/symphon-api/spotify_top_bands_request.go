package symphonapi

func (c *Client) GetTopBandsSpotify(bandChannel chan SpotifyResult) {
	topBands := []string{
		"Kings of Leon",
		"The Strokes",
		"The Killers",
		"The White Stripes",
		"The Black Keys",
		"Arctic Monkeys",
		"The Hives",
		"The Vines",
		"The Libertines",
		"The Raconteurs",
		"The Fratellis",
		"The Vaccines",
		"The Kooks",
		"The Wombats",
		"The Last Shadow Puppets",
	}
	bandChannel <- SpotifyResult{
		Message: topBands,
		Error:   nil,
	}
}
