package types

import "encoding/json"

type PositionDTO struct {
	X float32 `json:"x"`
	Z float32 `json:"z"`
}

func CreatePosition(bytes []byte) Position {
	pos := PositionDTO{}

	if err := json.Unmarshal(bytes, &pos); err != nil {
		panic(err)
	}

	return Position{
		x: pos.X,
		z: pos.Z,
	}
}

