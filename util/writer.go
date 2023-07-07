package util

import (
	"encoding/json"
	"io"
)

func ReadBody(r io.ReadCloser) (map[string]string, error) {
	var body map[string]string
	err := json.NewDecoder(r).Decode(&body)
	return body, err
}
