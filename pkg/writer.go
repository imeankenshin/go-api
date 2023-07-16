package pkg

import (
	"encoding/json"
	"io"
	"net/http"
)

func ReadBody(r io.ReadCloser) (map[string]string, error) {
	var body map[string]string
	err := json.NewDecoder(r).Decode(&body)
	return body, err
}

// Encode sets Content-Type header & writes the JSON encoding of v to w
func Encode(w http.ResponseWriter, v any) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return json.NewEncoder(w).Encode(v)
}
