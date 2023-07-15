package utils

import (
	"strings"
)


func LowercaseKeyMap(m map[string]string) map[string]string {
	// Create a new map to store the lowercase headers
	lowercaseKeyMap := make(map[string]string)

	// Convert the keys to lowercase and populate the lowercase headers map
	for key, value := range m {
		lowercaseKey := strings.ToLower(key)
		lowercaseKeyMap[lowercaseKey] = value
	}
	return lowercaseKeyMap
}