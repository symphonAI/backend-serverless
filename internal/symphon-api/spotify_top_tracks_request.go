package symphonapi

func (c *Client) GetTopTracksSpotify(bandChannel chan SpotifyResult) {
	topTracks := []string{
		"Use Somebody - Kings of Leon",
		"Last Nite - The Strokes",
		"Mr. Brightside - The Killers",
		"Seven Nation Army - The White Stripes",
		"Lonely Boy - The Black Keys",
		"I Bet You Look Good On The Dancefloor - Arctic Monkeys",
		"Hate To Say I Told You So - The Hives",
		"Get Free - The Vines",
		"Can't Stand Me Now - The Libertines",
		"Steady, As She Goes - The Raconteurs",
		"Chelsea Dagger - The Fratellis",
		"If You Wanna - The Vaccines",
		"Naive - The Kooks",
		"Let's Dance To Joy Division - The Wombats",
		"Standing Next To Me - The Last Shadow Puppets",
	}
	bandChannel <- SpotifyResult{
		Message: topTracks,
		Error:   nil,
	}
}
