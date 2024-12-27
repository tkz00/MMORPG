package utils

import (
	"encoding/json"
)

func ToJSON(s interface{}) ([]byte, error) {
	return json.Marshal(s)
}

func FromJSON(data []byte, s interface{}) error {
	return json.Unmarshal(data, &s)
}
