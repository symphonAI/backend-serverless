package symphonapi_test

import (
	symphonapi "backend-serverless/internal/symphon-api"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTracksWorksAsExpected(t *testing.T){
	m := symphonapi.ChatCompletionModel{}
	mockResponse := "{\"id\":\"chatcmpl-7pbvr9nUyd8CpbE8ZEeqoUtIsE378\",\"object\":\"chat.completion\",\"created\":1692534895,\"model\":\"gpt-3.5-turbo-0613\",\"choices\":[{\"index\":0,\"message\":{\"role\":\"assistant\",\"content\":\"{\\\"artist\\\": \\\"Cornelius\\\", \\\"track\\\": \\\"Mic Check\\\"}, \\n{\\\"artist\\\": \\\"Pizzicato Five\\\", \\\"track\\\": \\\"Twiggy Twiggy\\\"}, \\n{\\\"artist\\\": \\\"Stereolab\\\", \\\"track\\\": \\\"French Disko\\\"}, \\n{\\\"artist\\\": \\\"Fantastic Plastic Machine\\\", \\\"track\\\": \\\"Love is Psychedelic\\\"}, \\n{\\\"artist\\\": \\\"YMO\\\", \\\"track\\\": \\\"Firecracker\\\"}, \\n{\\\"artist\\\": \\\"Buffalo Daughter\\\", \\\"track\\\": \\\"Socks, Drugs and Rock'n'roll\\\"}, \\n{\\\"artist\\\": \\\"Plastics\\\", \\\"track\\\": \\\"Top Secret Man\\\"}, \\n{\\\"artist\\\": \\\"Suiyoubi no Campanella\\\", \\\"track\\\": \\\"Shakushain\\\"}, \\n{\\\"artist\\\": \\\"Towa Tei\\\", \\\"track\\\": \\\"Technova\\\"}, \\n{\\\"artist\\\": \\\"Takako Minekawa\\\", \\\"track\\\": \\\"Fantastic Cat\\\"}\"},\"finish_reason\":\"stop\"}],\"usage\":{\"prompt_tokens\":504,\"completion_tokens\":187,\"total_tokens\":691}}"
	tracks, err := m.ParseRecommendedTracksFromResponse([]byte(mockResponse))
	assert.NoError(t, err)
	assert.Equal(t, "Mic Check", tracks[0].Title)
	assert.Equal(t, "Cornelius", tracks[0].Artist)
}