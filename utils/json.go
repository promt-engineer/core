package utils

import "encoding/json"

func ReMarshal[T any](from interface{}) (T, error) {
	var to T

	b, err := json.Marshal(from)
	if err != nil {
		return to, err
	}

	return to, json.Unmarshal(b, &to)
}
