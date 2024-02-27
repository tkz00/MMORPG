package types

type GameDTO struct {
	Players map[string]PlayerDTO `json:"players"`
}
