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
	fmt.Println("Prompt:", prompt)
	return prompt, nil
}

func askForJSON(prompt string) string {
	ask := `You are a music recommendations AI. Please make music recommendations based on my prompt. Please respond with songs in the format {"artist": artist, "track": song}, within a JSON Array. Please return ` + NUMBER_OF_TRACKS + ` songs only. Please give me the response in JSON. Do not give me anything other than JSON. Here is my prompt, within tilde (~~~) signs: ~~~`
	return ask + prompt + "~~~"
}

func addTopBandsToPrompt(prompt string, topBands []string) string {
	prompt = prompt + " Within pound signs (£££) are my most-listened to musicians within the last month. Please use these to inform your music recommendations, but do not include music from these musicians: £££"
	for _, band := range topBands {
		prompt += band + ", "
	}
	return prompt + "£££"
}

func addTopTracksToPrompt(prompt string, topTracks []string) string {
	prompt = prompt + " Within caret signs (^^^) are my most-listened to musicians within the last month. Please use these to inform your music recommendations, but do not include these specific songs or musicians who performed these songs:  ^^^"
	for _, band := range topTracks {
		prompt += band + ", "
	}
	return prompt + "^^^"
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
