package promptengineering

const NUMBER_OF_TRACKS = "10"

func EngineerPrompt(prompt string, topBands []string, topTracks []string, options []string) (string, error) {
	prompt = askForJSON(prompt)


	return prompt, nil
}

func askForJSON(prompt string) string {
	ask := `. Please give me the data in the format {"artist": artist, "track": song}, within a JSON Array. Please give me ` + NUMBER_OF_TRACKS + `songs only. Please give me the response in JSON. Do not give me anything other than JSON. `;
	return ask + prompt 
}
