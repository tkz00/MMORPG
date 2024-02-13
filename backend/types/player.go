package types

import (
	"encoding/json"
)

type Player struct {
	position Position
}

func (p Player) CreatePlayer(x, z float32) Player {
	return Player{
		position: Position{X: x, Z: z},
	}
}
func CreatePosition(bytes []byte) Position{
	pos := Position{}
	if err := json.Unmarshal(bytes, &pos); err != nil {
        panic(err)
    }
	return pos
}
type Position struct {
	X float32 `json:"x"`
	Z float32 `json:"z"`
}