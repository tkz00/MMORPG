package types

import (
	"encoding/json"
)

type Position struct {
	X float32 `json:"x"`
	Z float32 `json:"z"`
}

func CreatePosition(bytes []byte) Position {
	pos := Position{}
	if err := json.Unmarshal(bytes, &pos); err != nil {
		panic(err)
	}
	return pos
}
