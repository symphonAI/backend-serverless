package main

import "encoding/json"

// TODO this file might have to go in a Lambda layer if
// we re-use it across a bunch of places

func structToJSON(data interface{}) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}