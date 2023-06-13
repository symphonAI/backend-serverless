package promptengineering

import "fmt"

const NUMBER_OF_TRACKS = "10"

func EngineerPrompt(prompt string, topBands []string, topTracks []string, options []string) (string, error) {
	prompt = prompt + ". "
	prompt = askForJSON(prompt)

	if !contains(options, "EXPLORE_MODE") {
		prompt = addTopBandsToPrompt(prompt, topBands)
		prompt = addTopTracksToPrompt(prompt, topTracks)
	} else {
		fmt.Println("EXPLORE_MODE SET")
	}
	return prompt, nil
}

func askForJSON(prompt string) string {
	ask := `You are a music recommendations AI. You make music recommendations based on the user prompt and their listening preferences. The user will prompt you for music. Please respond with data in the format {"artist": artist, "track": song}, within a JSON Array. Please give me ` + NUMBER_OF_TRACKS + ` songs only. Please give me the response in JSON. Do not give me anything other than JSON. Here is the user prompt: `
	return ask + prompt
}

func addTopBandsToPrompt(prompt string, topBands []string) string {
	prompt = prompt + " These are the user's most listened to bands from the last month: "
	for _, band := range topBands {
		prompt += band + ", "
	}
	return prompt
}

func addTopTracksToPrompt(prompt string, topTracks []string) string {
	prompt = prompt + " These are the user's most listened to songs from the last month: "
	for _, band := range topTracks {
		prompt += band + ", "
	}
	return prompt
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
